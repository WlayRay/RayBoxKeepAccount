package middleware

import (
	"ray_box/infrastructure/config"
	"ray_box/infrastructure/httputil"
	"ray_box/infrastructure/zlog"

	"github.com/kataras/iris/v12"
	"go.uber.org/zap"
)

func ParseToken(ctx iris.Context, Token string) {
	secret := config.GetConfig("SECRET_KEY")
	_, payload, err := httputil.VerifyToken(Token, secret)
	if err != nil {
		zlog.Error("token验证失败！", zap.Error(err))
	}

	ctx.Values().Set("uid", payload.UserDefined["uid"])
	ctx.Next()
}
