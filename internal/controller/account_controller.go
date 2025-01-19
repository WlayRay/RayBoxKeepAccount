package controller

import (
	"ray_box/infrastructure/httputil"
	"ray_box/infrastructure/xerror"
	"ray_box/infrastructure/zlog"
	"ray_box/internal/service"

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

	httputil.JSONSuccess().WithMsg("登录成功！").Response(ctx)
}
