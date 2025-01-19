package xerror

// 返回错误码模块
const (
	ErrCustom        = 10000 + iota // 自定义错误
	ErrRuntime                      // 运行错误
	ErrParamRequired                // 缺少参数
	ErrParamInvalid                 // 非法参数
	ErrForbidden                    // 没有权限
	ErrShow2User                    // 透传错误，提示语会直接展示在界面上给用户看

	ErrAddFail    // 创建失败
	ErrUpdateFail // 更新失败
	ErrDeleteFail // 删除失败
	ErrFindFail   // 获取失败
)

var baseErrorMap = errorDefinition{
	ErrRuntime:       "运行错误",
	ErrParamRequired: "缺少参数",
	ErrParamInvalid:  "非法参数",
	ErrForbidden:     " 没有权限",
	ErrAddFail:       " 创建失败",
	ErrUpdateFail:    " 更新失败",
	ErrDeleteFail:    " 删除失败",
	ErrFindFail:      " 获取失败",
}

type errorDefinition map[int]string

var allErrorMap = make(errorDefinition)

func init() {
	for k, v := range baseErrorMap {
		allErrorMap[k] = v
	}
}

// NewXErrorByCode 通过错误码构造错误
func NewXErrorByCode(code int) *XError {
	message := allErrorMap[code]

	return &XError{
		Code:    code,
		Message: message,
	}
}

// NewCustomXError 构造自定义错误
func NewCustomXError(message string) *XError {
	return &XError{
		Code:    ErrCustom,
		Message: message,
	}
}

// NewShow2UserXError 构造展示给用户的错误
func NewShow2UserXError(message string) *XError {
	return &XError{
		Code:    ErrShow2User,
		Message: message,
	}
}
