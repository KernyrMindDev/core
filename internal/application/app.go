package application

import (
	"net/http"
	"strings"

	emb "github.com/KernyrMindDev/core"
	"github.com/KernyrMindDev/core/internal/database"
	"github.com/KernyrMindDev/core/internal/handler"
	"github.com/KernyrMindDev/core/internal/model"
	"github.com/KernyrMindDev/core/internal/service"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

type App struct {
	Config        *model.AppConfig // 应用配置
	DB            *gorm.DB         // 数据库实例
	Engine        *gin.Engine      // Gin实例
	StaticFile    http.FileSystem  // 前端静态文件
	StaticHandler http.Handler     // 静态文件路由
}

func New(cfg *model.AppConfig) (*App, error) {
	// 初始化数据库
	db := database.InitDB(cfg)

	// 初始化JWT密钥
	service.InitJWTKey(db)

	// 切换Gin模式
	if !cfg.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	// 创建 Gin
	r := gin.Default()

	// 统一错误处理中间件，必须注册在所有业务路由之前，
	// 统一翻译 Handler/Service 通过 c.Error(err) 上报的错误
	r.Use(handler.ErrorHandlerMiddleware())

	// 注册路由
	api := r.Group("/api")
	// 初始化路由
	handler.InitRouter(api, db)

	// Swagger UI仅在调试模式下启动
	if cfg.Debug {
		// 自动重定向
		r.GET("/swagger", func(c *gin.Context) {
			c.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
		})
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	}

	app := &App{
		Config: cfg,
		DB:     db,
		Engine: r,
	}

	// 获取前端嵌入文件
	fileServer := emb.DistDirFS()
	app.StaticHandler = http.FileServer(fileServer)
	app.StaticFile = fileServer
	// 挂载
	r.NoRoute(app.StaticFileHandle)

	return app, nil
}

func (app *App) StaticFileHandle(c *gin.Context) {
	// 获取path
	path := c.Request.URL.Path
	// API 404
	if strings.HasPrefix(path, "/api/") {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not Found"})
		return
	}

	// 检查 embed 中是否存在该静态文件
	f, err := app.StaticFile.Open(strings.TrimPrefix(path, "/"))
	if err == nil {
		defer f.Close()
		app.StaticHandler.ServeHTTP(c.Writer, c.Request)
		return
	}

	// 不存在回退到 index.html
	c.Request.URL.Path = "/"
	app.StaticHandler.ServeHTTP(c.Writer, c.Request)
}
