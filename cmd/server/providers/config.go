// Package providers 提供依赖注入的 Provider 函数
package providers

import (
	"os"

	"github.com/gin-gonic/gin"

	"pets-server/internal/pkg/config"
)

// ProvideConfig 提供配置
func ProvideConfig() (*config.Config, error) {
	configPath := "configs/settings.yaml"
	if envPath := os.Getenv("CONFIG_PATH"); envPath != "" {
		configPath = envPath
	}

	cfg, err := config.Load(configPath)
	if err != nil {
		return nil, err
	}

	// 设置 Gin 模式
	gin.SetMode(cfg.Server.Mode)

	return cfg, nil
}
