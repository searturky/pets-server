// Package http HTTP 路由配置
package http

import (
	"github.com/gin-gonic/gin"

	"pets-server/internal/interfaces/http/handler"
	"pets-server/internal/interfaces/http/middleware"
)

const (
	// APIVersion API 版本号
	APIVersion = "v1"

	// APIBasePath API 基础路径
	APIBasePath = "/api"
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

	// 健康检查（不带版本号）
	router.GET("/health", healthCheck)

	// V1 API 路由组
	v1 := router.Group(APIBasePath + "/" + APIVersion)
	setupV1Routes(v1, cfg)

	return router
}

// healthCheck 健康检查
func healthCheck(c *gin.Context) {
	c.JSON(200, gin.H{
		"status":  "ok",
		"version": APIVersion,
	})
}

// setupV1Routes 配置 V1 版本路由
func setupV1Routes(api *gin.RouterGroup, cfg RouterConfig) {
	// 创建认证中间件
	authMiddleware := middleware.AuthMiddleware(
		middleware.JWTConfig{Secret: cfg.JWTSecret},
	)

	// 公开路由（无需认证）
	setupPublicRoutes(api, cfg)

	// 受保护路由（需要认证）
	setupProtectedRoutes(api, cfg, authMiddleware)
}

// setupPublicRoutes 配置公开路由（无需认证）
func setupPublicRoutes(api *gin.RouterGroup, cfg RouterConfig) {
	// 认证相关（登录、注册等）
	auth := api.Group("/auth")
	cfg.AuthHandler.RegisterRoutes(auth)
}

// setupProtectedRoutes 配置受保护路由（需要认证）
func setupProtectedRoutes(api *gin.RouterGroup, cfg RouterConfig, authMiddleware gin.HandlerFunc) {
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
