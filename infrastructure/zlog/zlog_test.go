package zlog_test

import (
	"errors"
	"ray_box/infrastructure/zlog"
	"testing"

	"go.uber.org/zap"
)

func TestZlog(t *testing.T) {
	failCount := 10086
	err := errors.New("test error")
	cacheKey := "test cache key"
	user := "test user"

	zlog.Info("GetCaptchaShow", zap.Any("failCount", failCount), zap.Error(err), zap.Any("cacheKey", cacheKey), zap.Any("user", user))

	zlog.Warn("GetCaptchaShow", zap.Any("failCount", failCount), zap.Error(err), zap.Any("cacheKey", cacheKey), zap.Any("user", user))

	zlog.Error("GetCaptchaShow", zap.Any("failCount", failCount), zap.Error(err), zap.Any("cacheKey", cacheKey), zap.Any("user", user))
}
