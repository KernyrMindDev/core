package dto

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
)

// 通用错误返回
type ErrorResponse struct {
	Error string `json:"error" example:"请求出错"`
	// Detail 仅用于开发调试，非 Debug 模式下不会出现在响应中
	Detail string `json:"detail,omitempty" example:"内部错误详情"`
}

// MarshalJSON 在非 Debug 模式下忽略 Detail 字段，
// 避免将内部错误细节暴露到生产环境响应中。
func (e ErrorResponse) MarshalJSON() ([]byte, error) {
	type Alias ErrorResponse

	if gin.Mode() == gin.DebugMode {
		return json.Marshal(Alias(e))
	}

	// 非 Debug 模式：强制清空 Detail 后再序列化
	alias := Alias(e)
	alias.Detail = ""
	return json.Marshal(alias)
}
