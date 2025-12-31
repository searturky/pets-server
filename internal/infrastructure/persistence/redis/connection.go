// Package redis Redis 连接和配置
package redis

import (
	"context"
	"fmt"
	"log"

	"github.com/go-redis/redis/v8"
)

// Config Redis 配置
type Config struct {
	Host     string
	Port     int
	Password string
	DB       int
	PoolSize int
}

// NewClient 创建 Redis 客户端
func NewClient(cfg Config) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
		PoolSize: cfg.PoolSize,
	})

	// 测试连接
	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect redis: %w", err)
	}

	log.Println("Redis connected successfully")
	return client, nil
}

