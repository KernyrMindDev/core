package config

import (
	"fmt"
	"os"

	"github.com/KernyrMindDev/core/internal/model"
	"gopkg.in/yaml.v3"
)

// 读取并解析应用配置
func ParseAppConfig(path string) (*model.AppConfig, error) {
	// 读取文件
	configFile, err := os.ReadFile(path)
	if err != nil {
		// 读取失败
		return nil, fmt.Errorf("Error: read config.yaml failed: %v", err)
	}
	var config model.AppConfig
	// 解析yaml
	err = yaml.Unmarshal(configFile, &config)
	if err != nil {
		// 解析出错
		return nil, err
	}
	return &config, nil
}
