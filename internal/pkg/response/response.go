// Package response HTTP 响应工具
package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 统一响应结构
type Response struct {
	Code    CustomCode `json:"code"`
	Message string     `json:"message"`
	Data    any        `json:"data,omitempty"`
}

// Success 成功响应
func Success(c *gin.Context, data any) {
	c.JSON(http.StatusOK, Response{
		Code:    CodeSuccess,
		Message: "success",
		Data:    data,
	})
}

// SuccessWithMessage 带消息的成功响应
func SuccessWithMessage(c *gin.Context, message string, data any) {
	c.JSON(http.StatusOK, Response{
		Code:    CodeSuccess,
		Message: message,
		Data:    data,
	})
}

// SuccessWithCode 带业务响应码的成功响应
func SuccessWithCode(c *gin.Context, code CustomCode, data any) {
	c.JSON(http.StatusOK, Response{
		Code:    code,
		Message: "success",
		Data:    data,
	})
}

// SuccessWithMessageAndCode 带业务响应码和消息的成功响应
func SuccessWithMessageAndCode(c *gin.Context, code CustomCode, message string, data any) {
	c.JSON(http.StatusOK, Response{
		Code:    code,
		Message: message,
		Data:    data,
	})
}

// Error 错误响应
func Error(c *gin.Context, code CustomCode, message string) {
	c.JSON(http.StatusInternalServerError, Response{
		Code:    code,
		Message: message,
	})
}

// ErrorWithData 带数据的错误响应
func ErrorWithData(c *gin.Context, code CustomCode, message string, data any) {
	c.JSON(http.StatusInternalServerError, Response{
		Code:    code,
		Message: message,
		Data:    data,
	})
}

// BadRequest 400 错误
func BadRequest(c *gin.Context, message string) {
	Error(c, CodeBadRequest, message)
}

// Unauthorized 401 错误
func Unauthorized(c *gin.Context, message string) {
	Error(c, CodeUnauthorized, message)
}

// Forbidden 403 错误
func Forbidden(c *gin.Context, message string) {
	Error(c, CodeForbidden, message)
}

// NotFound 404 错误
func NotFound(c *gin.Context, message string) {
	Error(c, CodeNotFound, message)
}

// InternalError 500 错误
func InternalError(c *gin.Context, message string) {
	Error(c, CodeInternalError, message)
}
