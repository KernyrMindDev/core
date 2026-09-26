package dto

// 通用错误返回
type ErrorResponse struct {
	Error string `json:"error" example:"请求出错"`
}
