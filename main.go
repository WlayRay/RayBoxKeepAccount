package main

import (
	_ "ray_box/infrastructure/config"
	"ray_box/infrastructure/httputil"
	"ray_box/infrastructure/xerror"
	"ray_box/infrastructure/zlog"
	"ray_box/internal/middleware"
	"ray_box/internal/route"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/pprof"
	"go.uber.org/zap"
)

func main() {
	defer zlog.Sync() // 确保在主函数退出前写入日志
	app := iris.New()

	// 处理内部服务器错误
	app.OnErrorCode(iris.StatusInternalServerError, func(ctx iris.Context) {
		httputil.JSONFailed().WithError(xerror.NewCustomXError("internal server error")).Response(ctx)
	})
	// 处理错误请求
	app.OnErrorCode(iris.StatusBadRequest, func(ctx iris.Context) {
		httputil.JSONFailed().WithError(xerror.NewCustomXError("bad request")).Response(ctx)
	})

	// 全局中间件，用于处理请求的记录和异常恢复
	app.UseGlobal(func(ctx iris.Context) {
		ctx.Record() // 记录请求信息
		defer func() {
			if err := recover(); err != nil { // 捕获异常
				zlog.Error("panic", zap.Any("error", err)) // 记录错误日志
				ctx.StatusCode(500)                        // 设置响应状态码为 500
			}
		}()
		ctx.Next()
	})

	app.UseGlobal(middleware.StartRequest)   // 注册请求开始的中间件
	app.DoneGlobal(middleware.FinishRequest) // 注册请求结束的中间件

	// 配置 pprof 性能分析路由
	p := pprof.New()
	app.Any("/debug/pprof", p)               // 主 pprof 路由
	app.Any("/debug/pprof/{action:path}", p) // 具体动作的 pprof 路由

	route.AccountRoute(app) // 账户相关

	// 启动 web 服务
	err := app.Listen(":8000")
	if err != nil {
		zlog.Error("listen error", zap.Error(err))
		panic(err)
	}
}
