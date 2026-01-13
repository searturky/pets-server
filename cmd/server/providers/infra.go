package providers

import (
	"fmt"
	"log"

	goredis "github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"pets-server/internal/domain/shared"
	"pets-server/internal/infrastructure/external/wechat"
	"pets-server/internal/infrastructure/messaging"
	"pets-server/internal/infrastructure/persistence/postgres"
	"pets-server/internal/infrastructure/persistence/redis"
	"pets-server/internal/pkg/config"
)

// ProvideDB 提供数据库连接
func ProvideDB(cfg *config.Config) (*gorm.DB, func(), error) {
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
		return nil, nil, fmt.Errorf("connect database: %w", err)
	}

	// 自动迁移
	if err := postgres.AutoMigrate(db); err != nil {
		return nil, nil, fmt.Errorf("migrate database: %w", err)
	}

	cleanup := func() {
		if sqlDB, err := db.DB(); err == nil {
			sqlDB.Close()
			log.Println("Database connection closed")
		}
	}

	return db, cleanup, nil
}

// ProvideRedis 提供 Redis 客户端
func ProvideRedis(cfg *config.Config) (*goredis.Client, func(), error) {
	client, err := redis.NewClient(redis.Config{
		Host:     cfg.Redis.Host,
		Port:     cfg.Redis.Port,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
		PoolSize: cfg.Redis.PoolSize,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("connect redis: %w", err)
	}

	cleanup := func() {
		client.Close()
		log.Println("Redis connection closed")
	}

	return client, cleanup, nil
}

// ProvideEventPublisher 提供事件发布器
func ProvideEventPublisher(cfg *config.Config) (shared.EventPublisher, func(), error) {
	// 优先使用 NATS
	if cfg.MQ.NATSURL != "" {
		publisher, err := messaging.NewNATSPublisher(messaging.Config{
			NATSURL:    cfg.MQ.NATSURL,
			StreamName: cfg.MQ.StreamName,
		})
		if err != nil {
			log.Printf("Warning: NATS unavailable (%v), trying Redis Stream", err)
		} else {
			return publisher, func() {
				publisher.Close()
				log.Println("NATS connection closed")
			}, nil
		}
	}

	// 尝试 Redis Stream 作为备用
	if cfg.Redis.Host != "" {
		publisher, err := messaging.NewRedisStreamPublisher(messaging.Config{
			RedisAddr:     fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
			RedisPassword: cfg.Redis.Password,
			RedisDB:       cfg.Redis.DB,
			StreamKey:     "game:events",
		})
		if err != nil {
			log.Printf("Warning: Redis Stream unavailable (%v), using noop publisher", err)
		} else {
			return publisher, func() {
				publisher.Close()
				log.Println("Redis Stream connection closed")
			}, nil
		}
	}

	// 使用 Noop 实现
	log.Println("MQ not configured, using noop publisher")
	return messaging.NewNoopPublisher(), func() {}, nil
}

// ProvideWechatAuth 提供微信认证服务
func ProvideWechatAuth(cfg *config.Config) *wechat.AuthService {
	return wechat.NewAuthService(wechat.Config{
		AppID:     cfg.Wechat.AppID,
		AppSecret: cfg.Wechat.AppSecret,
	})
}

// ProvideCacheService 提供缓存服务
func ProvideCacheService(client *goredis.Client) *redis.CacheService {
	return redis.NewCacheService(client)
}

// ProvideRankingStore 提供排行榜存储
func ProvideRankingStore(client *goredis.Client) *redis.RankingStore {
	return redis.NewRankingStore(client)
}

// ProvideUnitOfWork 提供工作单元
func ProvideUnitOfWork(db *gorm.DB) *postgres.UnitOfWork {
	return postgres.NewUnitOfWork(db)
}
