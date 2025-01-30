package controller

import (
	"crypto/rand"
	"io"
	"ray_box/infrastructure/config"
	"ray_box/infrastructure/httputil"
	"ray_box/infrastructure/xerror"
	"ray_box/infrastructure/zlog"
	"ray_box/internal/service"
	"time"

	"github.com/kataras/iris/v12"
	"go.uber.org/zap"
)

type accountController struct {
}

var (
	AccountController = accountController{}
)

func (*accountController) Login(ctx iris.Context) {
	type Params struct {
		Username string `json:"username" validate:"required"`
		Password string `json:"password" validate:"required"`
	}

	var params Params
	ok := httputil.GetRequestParams(ctx, &params)
	if !ok {
		httputil.JSONFailed().WithError(xerror.NewXErrorByCode(xerror.ErrParamInvalid)).Response(ctx)
		return
	}

	userInfo, err := service.AccountService.GetAccountInfo(params.Username)
	if err != nil {
		zlog.Error("获取用户信息失败！", zap.Error(err))
		httputil.JSONFailed().WithError(xerror.NewCustomXError("获取用户信息失败！")).Response(ctx)
		return
	}

	if userInfo["password"] != params.Password {
		zlog.Warn("用户名或密码错误！", zap.Any("userInfo", userInfo))
		httputil.JSONFailed().WithError(xerror.NewShow2UserXError("用户名或密码错误！")).Response(ctx)
		return
	}

	header := httputil.DefaultHeader
	payload := httputil.JwtPayload{
		Issue:      "RayBox",
		IssueAt:    time.Now().Unix(),
		Expiration: time.Now().Add(time.Hour * 24 * 15).Unix(),
		UserDefined: map[string]any{
			"username": params.Username,
		},
	}
	secret := config.GetConfig("SECRET_KEY")
	if token, err := httputil.GenerateToken(header, payload, secret); err != nil {
		zlog.Error("生成token失败！", zap.Error(err))
		httputil.JSONFailed().WithError(xerror.NewXErrorByCode(xerror.ErrRuntime)).Response(ctx)
	} else {
		refreshTokenBytes := make([]byte, 256)
		if _, err := io.ReadFull(rand.Reader, refreshTokenBytes); err != nil {
			zlog.Error("生成refresh token失败！", zap.Error(err))
			httputil.JSONFailed().WithError(xerror.NewXErrorByCode(xerror.ErrRuntime)).Response(ctx)
		}
		ctx.SetCookie(
			&iris.Cookie{
				Name:     "ray_box_token",
				Value:    string(refreshTokenBytes),
				Expires:  time.Now().Add(time.Hour * 24 * 30),
				Path:     "/",
				Domain:   "localhost", // 设置cookie的域名
				HttpOnly: false,
				Secure:   false,
			},
		)
		httputil.JSONSuccess().WithMsg("登录成功！").WithData(iris.Map{"Token": token}).Response(ctx)
	}
}
