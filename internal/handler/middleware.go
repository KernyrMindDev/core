package handler

import (
	"net/http"
	"strings"

	"github.com/KernyrMindDev/core/internal/dto"
	"github.com/KernyrMindDev/core/internal/service"
	"github.com/gin-gonic/gin"
)

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
			c.Abort() // 终止后续后续处理函数的执行
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
