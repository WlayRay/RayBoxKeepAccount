package xerror

import "fmt"

// XError 实现Error接口，用于定义接口统一返回的错误类型，包含错误码和错误信息，
// Code: 错误码，用于标识具体的错误类型，
// Message: 错误描述信息，部分可展示给用户。
type XError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *XError) Error() string {
	return fmt.Sprintf("code [%d] message [%s]", e.Code, e.Message)
}

func (e *XError) GetMessage() string {
	return e.Message
}

func (e *XError) GetCode() int {
	return e.Code
}
