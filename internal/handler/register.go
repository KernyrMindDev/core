// 集中实例化并注册路由
package handler

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitRouter(api *gin.RouterGroup, db *gorm.DB) {
	// 实例化路由
	authHandler := NewAuthHandler(db)
	boardHandler := NewBoardHandler(db)
	publicHandler := NewPublicHandler(db)

	// v1路由组
	v1 := api.Group("/v1")
	{
		v1.GET("/ping", publicHandler.PingHandler)
		// 用户认证
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
		}

		// 白板资源
		boards := v1.Group("/boards")
		// 加载Auth中间件
		boards.Use(JWTAuthMiddleware())
		// 路由组
		{
			boards.POST("", boardHandler.CreateBoard)
			boards.GET("/:id", boardHandler.GetBoardDetail)
			boards.GET("", boardHandler.GetBoards)
			boards.DELETE("/:id", boardHandler.DeleteBoard)
		}
	}
}
