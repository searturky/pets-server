// Package http HTTP 路由配置
//
// @title           宠物养成游戏 API
// @version         1.0
// @description     宠物养成游戏后端 REST API 文档
// @termsOfService  http://swagger.io/terms/
//
// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io
//
// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html
//
// @host      localhost:8080
// @BasePath  /api/v1
//
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
package http

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	authApp "pets-server/internal/application/auth"
	docs "pets-server/internal/interfaces/http/docs"
	"pets-server/internal/interfaces/http/handler"
	"pets-server/internal/interfaces/http/middleware"
	"pets-server/internal/pkg/config"
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
	SessionStore   authApp.SessionStore
	ServerMode     config.ServerMode // 服务器模式，用于控制 Swagger 开关
}

// NewRouter 创建路由
func NewRouter(cfg RouterConfig) *gin.Engine {
	router := gin.Default()

	// 全局中间件
	router.Use(middleware.CORSMiddleware())

	// 健康检查（不带版本号）
	router.GET("/health", healthCheck)

	// Swagger 文档路由（仅在开发环境：debug 模式）
	if cfg.ServerMode == config.ModeDebug {
		docs.SwaggerInfov1.Title = "宠物养成游戏 API 文档"
		docs.SwaggerInfov1.Description = "宠物养成游戏后端 REST API"
		docs.SwaggerInfov1.Version = "1.0"
		docs.SwaggerInfov1.BasePath = "/api/v1"
		docs.SwaggerInfov1.Schemes = []string{"http", "https"}

		router.GET("/swagger/v1/*any", ginSwagger.WrapHandler(
			swaggerFiles.Handler,
			ginSwagger.InstanceName("v1"),
		))
	}

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
		middleware.JWTConfig{
			Secret:       cfg.JWTSecret,
			SessionStore: cfg.SessionStore,
		},
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
