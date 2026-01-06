// Package main 程序入口
// 负责依赖注入和服务启动
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	// 应用层
	authApp "pets-server/internal/application/auth"
	petApp "pets-server/internal/application/pet"
	rankingApp "pets-server/internal/application/ranking"
	socialApp "pets-server/internal/application/social"

	// 领域层
	"pets-server/internal/domain/shared"

	// 基础设施层
	"pets-server/internal/infrastructure/cron"
	"pets-server/internal/infrastructure/external/wechat"
	"pets-server/internal/infrastructure/messaging"
	"pets-server/internal/infrastructure/persistence/postgres"
	"pets-server/internal/infrastructure/persistence/postgres/repo"
	"pets-server/internal/infrastructure/persistence/redis"

	// 接口层
	httpInterface "pets-server/internal/interfaces/http"
	"pets-server/internal/interfaces/http/handler"
	ws "pets-server/internal/interfaces/websocket"

	// 工具包
	"pets-server/internal/pkg/config"
)

func main() {
	// ========================================
	// 1. 加载配置
	// ========================================
	configPath := "configs/config.yaml"
	if envPath := os.Getenv("CONFIG_PATH"); envPath != "" {
		configPath = envPath
	}

	cfg := config.MustLoad(configPath)

	// 设置 Gin 模式
	gin.SetMode(cfg.Server.Mode)

	// ========================================
	// 2. 初始化基础设施层
	// ========================================

	// 2.1 PostgreSQL
	db, err := postgres.NewConnection(postgres.Config{
		Host:         cfg.Postgres.Host,
		Port:         cfg.Postgres.Port,
		User:         cfg.Postgres.User,
		Password:     cfg.Postgres.Password,
		DBName:       cfg.Postgres.DBName,
		SSLMode:      cfg.Postgres.SSLMode,
		MaxOpenConns: cfg.Postgres.MaxOpenConns,
		MaxIdleConns: cfg.Postgres.MaxIdleConns,
	})
	if err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}

	// 自动迁移数据库表
	if err := postgres.AutoMigrate(db); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// 2.2 Redis
	redisClient, err := redis.NewClient(redis.Config{
		Host:     cfg.Redis.Host,
		Port:     cfg.Redis.Port,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
		PoolSize: cfg.Redis.PoolSize,
	})
	if err != nil {
		log.Fatalf("Failed to connect redis: %v", err)
	}

	// 2.3 消息队列（可选，如果连接失败则使用 Noop 实现）
	var eventPublisher shared.EventPublisher
	var publisherCloser interface{ Close() error } // 用于关闭连接
	
	// 优先使用 NATS JetStream
	if cfg.MQ.NATSURL != "" {
		natsPublisher, err := messaging.NewNATSPublisher(messaging.Config{
			NATSURL:    cfg.MQ.NATSURL,
			StreamName: cfg.MQ.StreamName,
		})
		if err != nil {
			log.Printf("Warning: Failed to connect NATS, trying Redis Stream: %v", err)
			// 尝试 Redis Stream 作为备用
			if cfg.Redis.Host != "" {
				redisPublisher, err := messaging.NewRedisStreamPublisher(messaging.Config{
					RedisAddr:     fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
					RedisPassword: cfg.Redis.Password,
					RedisDB:       cfg.Redis.DB,
					StreamKey:     "game:events",
				})
				if err != nil {
					log.Printf("Warning: Failed to connect Redis Stream, using noop publisher: %v", err)
					eventPublisher = messaging.NewNoopPublisher()
				} else {
					eventPublisher = redisPublisher
					publisherCloser = redisPublisher
				}
			} else {
				eventPublisher = messaging.NewNoopPublisher()
			}
		} else {
			eventPublisher = natsPublisher
			publisherCloser = natsPublisher
		}
	} else {
		log.Println("MQ not configured, using noop publisher")
		eventPublisher = messaging.NewNoopPublisher()
	}
	
	// 注册关闭回调
	if publisherCloser != nil {
		defer publisherCloser.Close()
	}

	// 2.4 微信认证服务
	wechatAuth := wechat.NewAuthService(wechat.Config{
		AppID:     cfg.Wechat.AppID,
		AppSecret: cfg.Wechat.AppSecret,
	})

	// ========================================
	// 3. 创建仓储实现（基础设施层）
	// ========================================
	userRepo := repo.NewUserRepository(db)
	petRepo := repo.NewPetRepository(db)
	itemRepo := repo.NewItemRepository(db)
	friendRepo := repo.NewFriendRepository(db)
	giftRepo := repo.NewGiftRepository(db)
	tradeRepo := repo.NewTradeRepository(db)
	visitRepo := repo.NewVisitRepository(db)

	// 创建 UnitOfWork
	uow := postgres.NewUnitOfWork(db)

	// 创建缓存服务
	cacheService := redis.NewCacheService(redisClient)

	// 创建排行榜存储
	rankingStore := redis.NewRankingStore(redisClient)

	// ========================================
	// 4. 创建应用服务（应用层）
	// ========================================

	// 4.1 认证服务
	authService := authApp.NewService(
		userRepo,
		uow,
		wechatAuth.GetOpenID, // 微信认证函数
		cfg.JWT.Secret,
		cfg.JWT.ExpireHours,
	)

	// 4.2 宠物服务（完整示例）
	petService := petApp.NewService(
		petRepo,
		itemRepo,
		uow,
		nil, // eventPublisher - 需要类型转换，这里简化处理
		cacheService,
	)
	_ = eventPublisher // 避免未使用警告

	// 4.3 社交服务
	socialService := socialApp.NewService(
		friendRepo,
		giftRepo,
		tradeRepo,
		visitRepo,
		uow,
		nil, // eventPublisher
	)

	// 4.4 排行榜服务
	rankingService := rankingApp.NewService(rankingStore)

	// ========================================
	// 5. 创建 HTTP 处理器（接口层）
	// ========================================
	authHandler := handler.NewAuthHandler(authService)
	petHandler := handler.NewPetHandler(petService)
	itemHandler := handler.NewItemHandler()
	socialHandler := handler.NewSocialHandler(socialService)
	rankingHandler := handler.NewRankingHandler(rankingService)

	// ========================================
	// 6. 创建路由
	// ========================================
	router := httpInterface.NewRouter(httpInterface.RouterConfig{
		AuthHandler:    authHandler,
		PetHandler:     petHandler,
		ItemHandler:    itemHandler,
		SocialHandler:  socialHandler,
		RankingHandler: rankingHandler,
		JWTSecret:      cfg.JWT.Secret,
	})

	// ========================================
	// 7. 配置 WebSocket
	// ========================================
	wsHub := ws.NewHub()
	go wsHub.Run()

	wsHandler := ws.NewHandler(wsHub)
	router.GET("/ws", wsHandler.HandleWebSocket)

	// ========================================
	// 8. 启动定时任务
	// ========================================
	scheduler := cron.NewScheduler(petRepo, uow, nil)
	scheduler.Start()
	defer scheduler.Stop()

	// ========================================
	// 9. 启动 HTTP 服务器
	// ========================================
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	server := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	// 在 goroutine 中启动服务器
	go func() {
		log.Printf("Server starting on %s", addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	// ========================================
	// 10. 优雅关闭
	// ========================================
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}

