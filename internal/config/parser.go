package config

import (
	"fmt"
	"os"

	"github.com/KernyrMindDev/core/internal/model"
	"gopkg.in/yaml.v3"
)

// 读取并解析应用配置
func ParseAppConfig(path string) (*model.AppConfig, error) {
	configFile, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("Error: read config.yaml failed: %v", err)
	}
	var config model.AppConfig
	err = yaml.Unmarshal(configFile, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}
