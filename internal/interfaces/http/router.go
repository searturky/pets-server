// Package http HTTP 路由配置
package http

import (
	"github.com/gin-gonic/gin"

	"pets-server/internal/interfaces/http/handler"
	"pets-server/internal/interfaces/http/middleware"
)

// RouterConfig 路由配置
type RouterConfig struct {
	AuthHandler    *handler.AuthHandler
	PetHandler     *handler.PetHandler
	ItemHandler    *handler.ItemHandler
	SocialHandler  *handler.SocialHandler
	RankingHandler *handler.RankingHandler
	JWTSecret      string
}

// NewRouter 创建路由
func NewRouter(cfg RouterConfig) *gin.Engine {
	router := gin.Default()

	// 全局中间件
	router.Use(middleware.CORSMiddleware())

	// 健康检查
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// API 路由组
	api := router.Group("/api")
	{
		// 认证相关（无需登录）
		auth := api.Group("/auth")
		cfg.AuthHandler.RegisterRoutes(auth)

		// 需要登录的路由
		authMiddleware := middleware.AuthMiddleware(middleware.JWTConfig{Secret: cfg.JWTSecret})

		// 认证相关（需要登录）
		authProtected := api.Group("/auth")
		authProtected.Use(authMiddleware)
		cfg.AuthHandler.RegisterAuthRoutes(authProtected)

		// 宠物相关
		pet := api.Group("/pet")
		pet.Use(authMiddleware)
		cfg.PetHandler.RegisterRoutes(pet)

		// 道具相关
		items := api.Group("/items")
		items.Use(authMiddleware)
		cfg.ItemHandler.RegisterRoutes(items)

		// 社交相关
		social := api.Group("/social")
		social.Use(authMiddleware)
		cfg.SocialHandler.RegisterRoutes(social)

		// 排行榜
		ranking := api.Group("/ranking")
		ranking.Use(authMiddleware)
		cfg.RankingHandler.RegisterRoutes(ranking)
	}

	return router
}

