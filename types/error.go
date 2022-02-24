package types

// ErrorWithCode 支持自定义错误码的 error
type ErrorWithCode interface {
	error
	Code() int
}

// ErrorWithHTTPCode 支持自定义HTTP状态码的 error
type ErrorWithHTTPCode interface {
	error
	HTTPCode() int
}

// ErrorWithOrigin 支持返回原始错误的 error
type ErrorWithOrigin interface {
	error
	Origin() error
}

// CodeError 带错误码的错误
// 支持自定义错误码，HTTP状态码，及原始错误
// 注意：不要直接对 CodeError 进行比较操作，可能因为设置了不同的 origin 而不相等
type CodeError struct {
	httpStatus int
	errCode    int
	text       string
	origin     error
}

// NewCodeError 创建带错误码的错误
func NewCodeError(httpStatus int, errCode int, text string) CodeError {
	return CodeError{
		httpStatus: httpStatus,
		errCode:    errCode,
		text:       text,
	}
}

// Code 实现 envelope/ErrorWithCode 接口
func (e CodeError) Code() int {
	return e.errCode
}

// HTTPCode 实现 envelope/ErrorWithHTTPCode 接口
func (e CodeError) HTTPCode() int {
	return e.httpStatus
}

// Error 实现 error 接口
func (e CodeError) Error() string {
	return e.text
}

// Origin 实现 envelope/ErrorWithOrigin 接口
func (e CodeError) Origin() error {
	return e.origin
}

func (e CodeError) WithOrigin(err error) CodeError {
	dup := e // 复制一份
	dup.origin = err
	return dup
}
