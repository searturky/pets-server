// Package errors 错误定义
package errors

import "errors"

// 通用错误
var (
	ErrNotFound       = errors.New("资源不存在")
	ErrUnauthorized   = errors.New("未授权")
	ErrForbidden      = errors.New("无权限")
	ErrBadRequest     = errors.New("请求参数错误")
	ErrInternalServer = errors.New("服务器内部错误")
)

// AppError 应用错误
type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Err     error  `json:"-"`
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return e.Message + ": " + e.Err.Error()
	}
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

// New 创建应用错误
func New(code int, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
	}
}

// Wrap 包装错误
func Wrap(err error, code int, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// 预定义错误码
const (
	CodeSuccess          = 0
	CodeBadRequest       = 400
	CodeUnauthorized     = 401
	CodeForbidden        = 403
	CodeNotFound         = 404
	CodeConflict         = 409
	CodeInternalError    = 500
	CodeServiceUnavailable = 503
)

