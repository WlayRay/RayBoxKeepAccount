package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"gopkg.in/yaml.v3"
)

type Manager struct {
	configs sync.Map
}

var defaultManager = &Manager{}

const configDir = "config"

func init() {
	if err := defaultManager.Load(); err != nil {
		panic(fmt.Sprintf("加载配置文件失败: %v", err))
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
