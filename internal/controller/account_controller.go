package controller

import (
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"ray_box/common/tools"
	"ray_box/infrastructure/config"
	"ray_box/infrastructure/db"
	"ray_box/infrastructure/httputil"
	"ray_box/infrastructure/xerror"
	"ray_box/infrastructure/zlog"
	"ray_box/internal/service"
	"strconv"
	"sync"
	"time"

	"github.com/kataras/iris/v12"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type accountController struct {
}

var (
	AccountController = accountController{}
)

const (
	refreshCountKey = "refresh_token_used"
)

// Login 登录接口
func (*accountController) Login(ctx iris.Context) {
	type Params struct {
		Username string `json:"uid" validate:"required"`
		Password string `json:"pwd" validate:"required"`
	}

	var params Params
	httputil.GetRequestParams(ctx, &params)

	userInfo, err := service.AccountService.GetAccountInfo(params.Username)
	if errors.Is(err, gorm.ErrRecordNotFound) || userInfo["password"] != params.Password {
		httputil.JSONFailed().WithError(xerror.NewCustomXError("用户名或密码错误！")).Response(ctx)
		return
	} else if err != nil {
		zlog.Error("获取用户信息失败！", zap.Error(err))
		httputil.JSONFailed().WithError(xerror.NewXErrorByCode(xerror.ErrRuntime)).Response(ctx)
		return
	}

	// JWT的默认header
	header := httputil.DefaultHeader
	randomBytes := make([]byte, 512)
	if _, err := io.ReadFull(rand.Reader, randomBytes); err != nil {
		zlog.Error("生成随机字符串失败！", zap.Error(err))
		httputil.JSONFailed().WithError(xerror.NewXErrorByCode(xerror.ErrRuntime)).Response(ctx)
	}
	// payload
	payload := httputil.JwtPayload{
		Issue:      "RayBox",
		IssueAt:    time.Now().Unix(),
		Expiration: time.Now().Add(time.Hour * 24 * 7).Add(time.Hour * 24).Unix(),
		UserDefined: map[string]any{
			"uid":         params.Username,
			"randomBytes": randomBytes[:256],
		},
	}
	secret := config.GetConfig("SECRET_KEY")

	// 生成AccessToken
	if accessToken, err := httputil.GenerateToken(header, payload, secret); err != nil {
		zlog.Error("生成token失败！", zap.Error(err))
		httputil.JSONFailed().WithError(xerror.NewXErrorByCode(xerror.ErrRuntime)).Response(ctx)
	} else {
		// 生成RefreshToken
		payload.UserDefined["randomBytes"] = randomBytes[256:]
		refreshToken, _ := httputil.GenerateToken(header, payload, secret)
		ctx.SetCookie(
			&iris.Cookie{
				Name:     "ray_box_token",
				Value:    refreshToken,
				Expires:  time.Now().Add(time.Hour * 24 * 30),
				Path:     "/",
				Domain:   "localhost", // 设置cookie的域名
				HttpOnly: true,
				Secure:   false,
			},
		)

		timeoutContext, cancel := tools.GetTimeoutCtx(db.REDIS_DB_MASTER)
		defer cancel()

		redisConn := db.GetRedisConn(db.REDIS_DB_MASTER)
		key := fmt.Sprintf("%s:%s", refreshCountKey, params.Username)

		if setRedisErr := redisConn.Set(timeoutContext, key, 0, 0).Err(); setRedisErr == nil {
			zlog.Error("Redis调用失败！", zap.Error(setRedisErr))
			httputil.JSONFailed().WithError(xerror.NewXErrorByCode(xerror.ErrRuntime)).Response(ctx)
		}

		httputil.JSONSuccess().WithMsg("登录成功！").WithData(iris.Map{"Token": accessToken}).Response(ctx)
	}
}

// GetAuthToken 根据 AccessToken 获取 RefreshToken
func (*accountController) GetAuthToken(ctx iris.Context) {
	type Params struct {
		UID          string `json:"uid" validate:"required"`
		RefreshToken string `json:"token" validate:"required"`
	}

	var params Params
	httputil.GetRequestParams(ctx, &params)

	secret := config.GetConfig("SECRET_KEY")
	// 校验RefreshToken
	if _, payload, err := httputil.VerifyToken(params.RefreshToken, secret); err != nil {
		zlog.Error("token验证失败！", zap.Error(err))
	} else if payload.UserDefined["uid"] != params.UID {
		zlog.Error("token验证失败！", zap.Error(errors.New("身份校验失败！")))
	}

	timeoutContext, cancel := tools.GetTimeoutCtx(db.REDIS_DB_MASTER)
	defer cancel()

	redisConn := db.GetRedisConn(db.REDIS_DB_MASTER)
	key := fmt.Sprintf("%s:%s", refreshCountKey, params.UID)

	lock := sync.Mutex{}
	lock.Lock()
	countStr, _ := redisConn.Get(timeoutContext, key).Result()
	count, _ := strconv.Atoi(countStr)

	// 生成一个新的 AccessToken
	// 默认的 header
	header := httputil.DefaultHeader
	randomBytes := make([]byte, 256)
	if _, err := io.ReadFull(rand.Reader, randomBytes); err != nil {
		zlog.Error("生成随机字符串失败！", zap.Error(err))
		httputil.JSONFailed().WithError(xerror.NewXErrorByCode(xerror.ErrRuntime)).Response(ctx)
	}
	// payload
	payload := httputil.JwtPayload{
		Issue:      "RayBox",
		IssueAt:    time.Now().Unix(),
		Expiration: time.Now().Add(time.Hour * 24 * 7).Add(time.Hour * 24).Unix(),
		UserDefined: map[string]any{
			"uid":         params.UID,
			"randomBytes": randomBytes,
		},
	}
	newAccessToken, _ := httputil.GenerateToken(header, payload, secret)
	result := iris.Map{
		"AccessToken": newAccessToken,
	}

	//每个 RefreshToken 只能使用20次
	if count > 20 {
		if _, err := io.ReadFull(rand.Reader, randomBytes); err != nil {
			zlog.Error("生成随机字符串失败！", zap.Error(err))
			httputil.JSONFailed().WithError(xerror.NewXErrorByCode(xerror.ErrRuntime)).Response(ctx)
		}
		newRefreshToken, _ := httputil.GenerateToken(header, payload, secret)
		ctx.SetCookie(
			&iris.Cookie{
				Name:     "ray_box_token",
				Value:    newRefreshToken,
				Expires:  time.Now().Add(time.Hour * 24 * 30),
				Path:     "/",
				Domain:   "localhost", // 设置 cookie 的域名
				HttpOnly: true,
				Secure:   false,
			},
		)

		if setRedisErr := redisConn.Set(timeoutContext, key, 0, 0).Err(); setRedisErr == nil {
			zlog.Error("Redis调用失败！", zap.Error(setRedisErr))
			httputil.JSONFailed().WithError(xerror.NewXErrorByCode(xerror.ErrRuntime)).Response(ctx)
		}
	}
	lock.Unlock()

	httputil.JSONSuccess().WithData(result).Response(ctx)

}
