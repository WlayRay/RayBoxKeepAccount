package httputil

import (
	"ray_box/infrastructure/zlog"

	"github.com/bytedance/sonic"
	"github.com/kataras/iris/v12"
	"go.uber.org/zap"
)

// 从已解析的结果中，提取请求体内容并转换为字节数组返回
func GetRequestBodyFromValues(ctx iris.Context) []byte {
	bodyBytes, err := ctx.GetBody()
	if err != nil {
		zlog.Error("get request body error", zap.Error(err))
		panic(err)
	}
	return bodyBytes
}

// GetRequestParams 获取参数, 1.解析json; 2.校验validate
func GetRequestParams(ctx iris.Context, params any) bool {
	bodyBytes := GetRequestBodyFromValues(ctx)
	if len(bodyBytes) == 0 {
		zlog.Warn("empty request body")
	}

	if jErr := sonic.Unmarshal(bodyBytes, &params); jErr != nil {
		zlog.Error("json unmarshal error", zap.Error(jErr))
		return false
	}

	v := GetValidator()
	if vErr := v.Struct(params); vErr != nil {
		zlog.Error("validate error", zap.Error(vErr))
		return false
	}

	return true
}
