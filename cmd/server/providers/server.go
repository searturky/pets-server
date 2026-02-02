package providers

import (
	"github.com/gin-gonic/gin"

	authApp "pets-server/internal/application/auth"
	"pets-server/internal/infrastructure/cron"
	"pets-server/internal/infrastructure/persistence/postgres"
	httpInterface "pets-server/internal/interfaces/http"
	"pets-server/internal/interfaces/http/handler"
	ws "pets-server/internal/interfaces/websocket"
	"pets-server/internal/pkg/config"
)

// ProvideRouter 提供路由
func ProvideRouter(cfg *config.Config, services *ServiceSet, wsHandler *ws.Handler, sessionStore authApp.SessionStore) *gin.Engine {
	// 创建 HTTP Handler
	authHandler := handler.NewAuthHandler(services.Auth)
	petHandler := handler.NewPetHandler(services.Pet)
	itemHandler := handler.NewItemHandler()
	socialHandler := handler.NewSocialHandler(services.Social)
	rankingHandler := handler.NewRankingHandler(services.Ranking)

	// 创建路由
	router := httpInterface.NewRouter(httpInterface.RouterConfig{
		AuthHandler:    authHandler,
		PetHandler:     petHandler,
		ItemHandler:    itemHandler,
		SocialHandler:  socialHandler,
		RankingHandler: rankingHandler,
		JWTSecret:      cfg.JWT.Secret,
		SessionStore:   sessionStore,
		ServerMode:     cfg.Server.Mode, // 传递服务器模式
	})

	// 注册 WebSocket 路由
	router.GET("/ws", wsHandler.HandleWebSocket)

	return router
}

// ProvideWSHub 提供 WebSocket Hub
func ProvideWSHub() *ws.Hub {
	return ws.NewHub()
}

// ProvideWSHandler 提供 WebSocket Handler
func ProvideWSHandler(hub *ws.Hub) *ws.Handler {
	return ws.NewHandler(hub)
}

// ProvideScheduler 提供定时任务调度器
func ProvideScheduler(repos *RepoSet, uow *postgres.UnitOfWork) *cron.Scheduler {
	return cron.NewScheduler(repos.Pet, uow, nil)
}
