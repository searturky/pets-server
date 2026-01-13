//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"

	"pets-server/cmd/server/providers"
)

// InitializeApp 使用 Wire 初始化应用
// 这个函数的实现由 Wire 工具自动生成
func InitializeApp() (*App, func(), error) {
	wire.Build(
		// 配置
		providers.ProvideConfig,

		// 基础设施
		providers.ProvideDB,
		providers.ProvideRedis,
		providers.ProvideEventPublisher,
		providers.ProvideWechatAuth,
		providers.ProvideCacheService,
		providers.ProvideRankingStore,
		providers.ProvideUnitOfWork,

		// 仓储
		providers.ProvideRepoSet,

		// 服务
		providers.ProvideServiceSet,

		// HTTP & WebSocket
		providers.ProvideWSHub,
		providers.ProvideWSHandler,
		providers.ProvideRouter,
		providers.ProvideScheduler,

		// App
		NewApp,
	)

	return nil, nil, nil
}
