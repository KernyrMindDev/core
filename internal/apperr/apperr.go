package apperr

import (
	"errors"
	"net/http"
)

// AppError 携带 HTTP 状态码和对外展示信息的业务错误
type AppError struct {
	HTTPStatus int    // HTTP状态码
	Message    string // 对外展示的错误消息
	Err        error  // 原始错误
}

func (e *AppError) Error() string {
	return e.Message
}

// Unwrap 使 errors.Is / errors.As 能够穿透到底层原始错误
func (e *AppError) Unwrap() error {
	return e.Err
}

// NotFound 资源不存在 (404)
func NotFound(message string) *AppError {
	return &AppError{HTTPStatus: http.StatusNotFound, Message: message}
}

// Forbidden 已登录但权限不足 (403)
func Forbidden(message string) *AppError {
	return &AppError{HTTPStatus: http.StatusForbidden, Message: message}
}

// Unauthorized 未登录或凭证无效 (401)
func Unauthorized(message string) *AppError {
	return &AppError{HTTPStatus: http.StatusUnauthorized, Message: message}
}

// BadRequest 请求参数错误 (400)
func BadRequest(message string) *AppError {
	return &AppError{HTTPStatus: http.StatusBadRequest, Message: message}
}

// Internal 服务器内部错误 (500)。err 为原始错误
func Internal(message string, err error) *AppError {
	return &AppError{HTTPStatus: http.StatusInternalServerError, Message: message, Err: err}
}

// As 尝试将 err 转换为 *AppError，供错误处理中间件判断
func As(err error) (*AppError, bool) {
	var appErr *AppError
	ok := errors.As(err, &appErr)
	return appErr, ok
}
