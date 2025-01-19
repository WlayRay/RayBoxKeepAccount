package middleware

import (
	"os"
	"ray_box/infrastructure/zlog"
	"strconv"
	"time"

	"github.com/kataras/iris/v12"
	"go.uber.org/zap"
)

func StartRequest(ctx iris.Context) {
	ctx.Values().Set("request_start_time", time.Now().UnixNano())
	ctx.Next()
}

func FinishRequest(ctx iris.Context) {
	CheckReqProcessTime(ctx)
	ctx.Next()
}

func CheckReqProcessTime(ctx iris.Context) {
	requestStart := ctx.Values().Get("request_start_time")
	if start, ok := requestStart.(int64); ok {
		now := time.Now().UnixNano()
		cost := (now - start) / int64(time.Millisecond)

		var numThreshold int64 = 2000
		slowRequestThreshold, envOK := os.LookupEnv("SLOW_REQUEST_THRESHOLD")
		if num, err := strconv.Atoi(slowRequestThreshold); envOK && err == nil {
			numThreshold = int64(num)
		}

		if cost > numThreshold {
			zlog.Warn("processing request too slow",
				zap.String("uri", ctx.FullRequestURI()),
				zap.String("cost", strconv.FormatInt(cost, 10)+"ms"),
			)
		}
	} else {
		zlog.Error("context value request_start type error")
	}
}
