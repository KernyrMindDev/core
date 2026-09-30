package handler

import (
	"net/http"
	"strings"

	"github.com/KernyrMindDev/core/internal/apperr"
	"github.com/KernyrMindDev/core/internal/dto"
	"github.com/KernyrMindDev/core/internal/service"
	"github.com/gin-gonic/gin"
)

// ErrorHandlerMiddleware 统一错误处理中间件
//
// 统一读取 c.Errors 并翻译成 HTTP 响应
//
// 必须注册在路由链最前面
//
// 业务 Handler/Service 只需要 c.Error(err) 上报错误并 return，
func ErrorHandlerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 {
			return
		}

		err := c.Errors.Last().Err
		if appErr, ok := apperr.As(err); ok {
			c.JSON(appErr.HTTPStatus, dto.ErrorResponse{
				Error:  appErr.Message,
				Detail: errString(appErr.Err),
			})
			return
		}

		// 未知错误，统一按 500 处理，避免泄露内部细节
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:  "服务器内部错误",
			Detail: err.Error(),
		})
	}
}

func errString(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}

// JWTAuthMiddleware JWT认证中间件
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取 Authorization Header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
				Error:  "Unauthorized",
				Detail: "Lost Authorization header",
			})
			c.Abort() // 终止后续处理函数的执行
			return
		}

		// 解析 Bearer Token 格式
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
				Error:  "Unauthorized",
				Detail: "Expect Bearer token",
			})
			c.Abort()
			return
		}

		// 验证 Token
		claims, err := service.ParseToken(parts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
				Error:  "Unauthorized",
				Detail: "Invalid or expired token",
			})
			c.Abort()
			return
		}

		// 将解析出的用户信息保存到全局上下文 Context 中
		c.Set("uid", claims.Uid)

		c.Next() // 执行后续的路由处理函数
	}
}
