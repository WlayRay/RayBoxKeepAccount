package zlog

import (
	"errors"
	"strings"
	"syscall"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

/*
参考 log 标准库对 zap log 进行封装
日志输出基础内容：
caller： 调用者(某个文件某一行报错)
stacktrace: 错误堆栈(暂不启动)
level： 日志级别
ts： 调用日志的时间戳
*/
var zLog = GetLogger()

func GetLogger() *zap.Logger {
	config := zap.NewProductionConfig()
	config.Level = zap.NewAtomicLevelAt(GetCurrentLogLevel())
	newLogger, err := config.Build(
		// zap.AddStacktrace(zap.ErrorLevel),
		zap.AddCaller(),
		zap.AddCallerSkip(1),
	)
	if err != nil {
		panic("Failed to build logger: " + err.Error())
	}
	return newLogger
}

// GetCurrentLogLevel 获取当前的日志级别
func GetCurrentLogLevel() zapcore.Level {
	const defaultLevel = zap.WarnLevel
	// envLogLevel, _ := os.LookupEnv("LIB_LOG_LEVEL")
	envLogLevel := "info"
	var level zapcore.Level
	switch strings.ToLower(envLogLevel) {
	case "debug":
		level = zap.DebugLevel
	case "info":
		level = zap.InfoLevel
	case "warning":
		level = zap.WarnLevel
	case "error":
		level = zap.ErrorLevel
	case "dpanic":
		level = zap.DPanicLevel
	case "panic":
		level = zap.PanicLevel
	case "fatal":
		level = zap.FatalLevel
	default:
		level = defaultLevel
	}

	return level
}

// Error 封装 zap log Error 方法
func Error(msg string, fields ...zap.Field) {
	zLog.Error(msg, fields...)
}

// Warn 封装 zap log Warn 方法
func Warn(msg string, fields ...zap.Field) {
	zLog.Warn(msg, fields...)
}

// Info 封装 zap log Info 方法
func Info(msg string, fields ...zap.Field) {
	zLog.Info(msg, fields...)
}

// Debug 封装 zap log Debug 方法
func Debug(msg string, fields ...zap.Field) {
	zLog.Debug(msg, fields...)
}

// Panic 封装 zap log Panic 方法
func Panic(msg string, fields ...zap.Field) {
	zLog.Panic(msg, fields...)
}

// Sync 调用底层 Sync， main退出前调用一次即可
// https://github.com/uber-go/zap/issues/1093
func Sync() {
	err := zLog.Sync()
	if err != nil && !errors.Is(err, syscall.ENOTTY) && !errors.Is(err, syscall.EINVAL) {
		zLog.Error("zLog Sync", zap.Any("err", err))
		return
	}
}
