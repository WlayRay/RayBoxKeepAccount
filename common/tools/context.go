package tools

import (
	"context"
	"ray_box/infrastructure/config"
	"strconv"
	"strings"
	"time"
)

var (
	timeoutContext          context.Context
	cancelFunc              context.CancelFunc
	RedisTimeoutDuration    = config.GetConfig("REDIS_TIMEOUT_DURATION")
	PostgresTimeoutDuration = config.GetConfig("POSTGRES_TIMEOUT_DURATION")
)

func GetTimeoutCtx(dbName string) (context.Context, context.CancelFunc) {
	var (
		durationStr string
		duration    time.Duration
	)
	dbName = strings.ToLower(dbName)
	switch dbName {
	case "redis":
		durationStr = config.GetConfig("REDIS_TIMEOUT_DURATION")
	case "postgres":
		durationStr = config.GetConfig("POSTGRES_TIMEOUT_DURATION")
	default:
		panic("dbName错误")
	}

	durationInt, err := strconv.Atoi(durationStr)
	if err != nil {
		panic(err)
	}

	duration = time.Duration(time.Second * time.Duration(durationInt))
	timeoutContext, cancelFunc = context.WithTimeout(context.Background(), duration)

	return timeoutContext, cancelFunc
}
