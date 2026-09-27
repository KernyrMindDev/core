package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"

	emb "github.com/KernyrMindDev/core"
	"github.com/KernyrMindDev/core/internal/config"
	"github.com/KernyrMindDev/core/internal/database"
	"github.com/KernyrMindDev/core/internal/handler"
	"github.com/KernyrMindDev/core/internal/service"

	_ "github.com/KernyrMindDev/core/docs"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/gin-gonic/gin"
)

// @title           KernyrMind Backend
// @version         1.0
// @description     API Docs for KernyrMind Backend
// @host            localhost:8080
// @BasePath        /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description 在输入框中填入: Bearer <Token>

func main() {
	// 解析config.yaml命令行参数
	var configPath string
	flag.StringVar(&configPath, "c", "./config.yaml", "path to configuration file")
	flag.Parse()

	// 检查文件是否存在
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		fmt.Printf("Error: config file not found at '%s'\n", configPath)
		os.Exit(1)
	}

	// 解析yaml
	appconfig, err := config.PasreAppConfig(configPath)
	if err != nil {
		fmt.Printf("Error: config file pasre failed: %v\n", err)
		os.Exit(1)
	}

	// 初始化数据库
	db := database.InitDB(appconfig)

	// 生产模式切换
	if !appconfig.Debug {
		// fmt.Println("Switch to release mode")
		gin.SetMode(gin.ReleaseMode)
	}

	// 初始化JWT密钥
	service.InitJWTKey(db)

	// 初始化Gin路由
	r := gin.Default()

	// 实例化Handler并注入DB
	authHandler := handler.NewAuthHandler(db)
	boardHandler := handler.NewBoardHandler(db)
	publicHandler := handler.NewPublicHandler(db)

	// 注册路由
	api := r.Group("/api/v1")
	{
		api.GET("/ping", publicHandler.PingHandler)

		// 用户认证
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
		}

		// 白板资源
		boards := api.Group("/boards")
		// 加载Auth中间件
		boards.Use(handler.JWTAuthMiddleware())
		// 路由组
		{
			boards.POST("", boardHandler.CreateBoard)
			boards.GET("/:id", boardHandler.GetBoardDetail)
		}
	}

	// Swagger UI仅在调试模式下启动
	if appconfig.Debug {
		// 自动重定向
		r.GET("/swagger", func(c *gin.Context) {
			c.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
		})
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	}

	// 获取前端嵌入文件
	fileServer := emb.DistDirFS()
	staticHandler := http.FileServer(fileServer)
	r.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path

		// API 404
		if strings.HasPrefix(path, "/api/") {
			c.JSON(http.StatusNotFound, gin.H{"error": "Not Found"})
			return
		}

		// 检查 embed 中是否存在该静态文件
		f, err := fileServer.Open(strings.TrimPrefix(path, "/"))
		if err == nil {
			defer f.Close()
			staticHandler.ServeHTTP(c.Writer, c.Request)
			return
		}

		// 不存在回退到 index.html
		c.Request.URL.Path = "/"
		staticHandler.ServeHTTP(c.Writer, c.Request)
	})

	r.Run(fmt.Sprintf(":%v", appconfig.Port))
}
