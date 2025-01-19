package httputil

import (
	"ray_box/infrastructure/xerror"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
)

// response 响应结构
type response struct {
	HTTPCode    int
	HTTPHeader  map[string]string
	BodySuccess bool
	BodyMsg     string
	BodyCode    int
	BodyData    any
}

// successResponse 成功的响应结构
type successResponse struct {
	response
}

// failedResponse 失败的响应结构
type failedResponse struct {
	response
}

// JSONSuccess 成功的响应，
// BodyCode 默认为 200。
func JSONSuccess() *successResponse {
	return &successResponse{
		response: response{
			HTTPCode:    iris.StatusOK,
			BodyCode:    200,
			BodySuccess: true,
		},
	}
}

// JSONFailed 失败的响应
func JSONFailed() *failedResponse {
	return &failedResponse{
		response: response{
			HTTPCode:    iris.StatusOK,
			BodySuccess: false,
		},
	}
}

func (r *successResponse) WithData(data any) *successResponse {
	r.BodyData = data
	return r
}

func (r *successResponse) WithMsg(msg string) *successResponse {
	r.BodyMsg = msg
	return r
}

func (r *response) WithHeader(header map[string]string) *response {
	r.HTTPHeader = header
	return r
}

func (r *failedResponse) WithError(err error) *failedResponse {
	if xErr, ok := err.(*xerror.XError); ok {
		var msg string
		if len(xErr.Message) > 0 {
			msg = xErr.Message
		}
		r.BodyCode = xErr.Code
		r.BodyMsg = msg
	} else {
		r.HTTPCode = iris.StatusInternalServerError
		r.BodyCode = 500
	}
	return r
}

func (r *response) Response(ctx iris.Context) {
	ctx.StatusCode(r.HTTPCode)
	for k, v := range r.HTTPHeader {
		ctx.Header(k, v)
	}
	err := ctx.JSON(iris.Map{
		"success": r.BodySuccess,
		"data":    r.BodyData,
		"msg":     r.BodyMsg,
		"code":    r.BodyCode,
	}, context.DefaultJSONOptions)
	if err != nil {
		panic(err)
	}
}
