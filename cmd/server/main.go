package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/KernyrMindDev/core/internal/application"
	"github.com/KernyrMindDev/core/internal/config"

	_ "github.com/KernyrMindDev/core/docs"
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

	// 解析yaml
	appconfig, err := config.ParseAppConfig(configPath)
	if err != nil {
		fmt.Printf("Error: config file pasre failed: %v\n", err)
		os.Exit(1)
	}

	app, err := application.New(appconfig)
	if err != nil {
		fmt.Printf("Failed when start: %v", err)
	}

	// 启动应用
	app.Engine.Run(fmt.Sprintf(":%v", appconfig.Port))
}
