package config

import (
	"fmt"
	"os"
	"path/filepath"
	"ray_box/infrastructure/zlog"
	"runtime"
	"strings"
	"sync"

	"github.com/redis/go-redis/v9"
	"gopkg.in/yaml.v3"
)

type Manager struct {
	configs sync.Map
}

var (
	defaultManager = &Manager{}
	rootPath       string
	configDir      string
)

func init() {
	_, filename, _, _ := runtime.Caller(0)
	rootPath = filepath.Join(filepath.Dir(filename), "../..")
	configDir = filepath.Join(rootPath, "config")

	if err := defaultManager.Load(); err != nil {
		panic(fmt.Sprintf("加载配置文件失败 -> %v", err))
	}
}

// Load 加载所有配置文件
func (m *Manager) Load() error {
	entries, err := os.ReadDir(configDir)
	if err != nil {
		return fmt.Errorf("读取配置目录失败: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		if !strings.HasSuffix(name, ".yaml") && !strings.HasSuffix(name, ".yml") {
			continue
		}

		if err := m.loadYAMLFile(filepath.Join(configDir, name)); err != nil {
			return err
		}
	}
	return nil
}

// loadYAMLFile 加载单个YAML文件
func (m *Manager) loadYAMLFile(filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("读取文件失败: %w", err)
	}

	config := make(map[string]string)
	if err := yaml.Unmarshal(data, &config); err != nil {
		return fmt.Errorf("解析YAML失败: %w", err)
	}

	for key, value := range config {
		m.configs.Store(key, m.resolveValue(value))
	}

	return nil
}

// resolveValue 处理配置值中的环境变量
func (m *Manager) resolveValue(value string) string {
	parts := strings.Split(value, "|")
	if len(parts) != 2 {
		return value
	}

	envKey := strings.Trim(parts[0], "${}")
	envValue := os.Getenv(envKey)
	if envValue != "" {
		return envValue
	}
	return parts[1]
}

// GetConfig 获取配置值
func GetConfig(key string) string {
	return defaultManager.Get(key)
}

// Get 获取配置值的方法实现
func (m *Manager) Get(key string) string {
	value, ok := m.configs.Load(key)
	if !ok {
		return ""
	}
	return value.(string)
}

// GetWithDefault 获取配置值，如果不存在返回默认值
func GetWithDefault(key, defaultValue string) string {
	value := defaultManager.Get(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// RedisConfig 获取redis配置
func RedisConfig(name string) *redis.UniversalOptions {
	name = strings.ToUpper(name)
	var option redis.UniversalOptions
	// 当UniversalOptions Address > 1时为集群
	if val, ok := defaultManager.configs.Load(fmt.Sprintf("REDIS_%s_NODES", name)); ok && val != "" {
		addrs := strings.Split(val.(string), ",")
		option = redis.UniversalOptions{
			Addrs: addrs,
			// Password:      defaultManager.Get(fmt.Sprintf("REDIS_%s_PASSWORD", name)),
			PoolSize:     60,
			MinIdleConns: 30,
		}
	} else {
		zlog.Warn(fmt.Sprintf("REDIS_%s_NODES", name) + "不存在")
		_, ok = defaultManager.configs.Load(fmt.Sprintf("REDIS_%s_HOST", name))
		if !ok {
			zlog.Warn(fmt.Sprintf("REDIS_%s_HOST", name) + "不存在")
		}
		option = redis.UniversalOptions{
			Addrs:    []string{defaultManager.Get(fmt.Sprintf("REDIS_%s_HOST", name)) + ":" + defaultManager.Get(fmt.Sprintf("REDIS_%s_PORT", name))},
			Password: defaultManager.Get(fmt.Sprintf("REDIS_%s_PASSWORD", name)),
			DB:       0,
			PoolSize: 100,
		}
	}

	return &option
}
