// Package messaging 消息队列事件发布实现
package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/redis/go-redis/v9"

	"pets-server/internal/domain/shared"
)

// Config MQ 配置
type Config struct {
	// NATS 配置
	NATSURL    string // nats://localhost:4222
	StreamName string // 默认: game-events

	// Redis 配置（备用方案）
	RedisAddr     string // localhost:6379
	RedisPassword string
	RedisDB       int
	StreamKey     string // 默认: game:events
}

// --- NATS JetStream 实现（主要方案） ---

// NATSPublisher NATS JetStream 事件发布器
type NATSPublisher struct {
	nc         *nats.Conn
	js         nats.JetStreamContext
	streamName string
}

// NewNATSPublisher 创建 NATS JetStream 发布器
// 依赖: go get github.com/nats-io/nats.go
func NewNATSPublisher(cfg Config) (*NATSPublisher, error) {
	// 连接 NATS 服务器
	nc, err := nats.Connect(cfg.NATSURL,
		nats.MaxReconnects(-1),          // 无限重连
		nats.ReconnectWait(time.Second), // 重连间隔
		nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
			if err != nil {
				log.Printf("NATS disconnected: %v", err)
			}
		}),
		nats.ReconnectHandler(func(nc *nats.Conn) {
			log.Printf("NATS reconnected to %s", nc.ConnectedUrl())
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to NATS: %w", err)
	}

	// 创建 JetStream 上下文
	js, err := nc.JetStream()
	if err != nil {
		nc.Close()
		return nil, fmt.Errorf("failed to create JetStream context: %w", err)
	}

	streamName := cfg.StreamName
	if streamName == "" {
		streamName = "game-events"
	}

	// 创建或更新 Stream（用于持久化事件）
	_, err = js.AddStream(&nats.StreamConfig{
		Name:      streamName,
		Subjects:  []string{"game.>"},      // 监听所有 game.* 主题
		Retention: nats.LimitsPolicy,       // 按大小/数量/时间限制保留
		MaxAge:    7 * 24 * time.Hour,      // 保留7天
		MaxBytes:  10 * 1024 * 1024 * 1024, // 最大10GB
		Storage:   nats.FileStorage,        // 文件存储（持久化）
		Replicas:  1,                       // 副本数（单机为1）
	})
	if err != nil && err != nats.ErrStreamNameAlreadyInUse {
		nc.Close()
		return nil, fmt.Errorf("failed to create stream: %w", err)
	}

	log.Printf("NATS JetStream connected successfully (stream: %s)", streamName)

	return &NATSPublisher{
		nc:         nc,
		js:         js,
		streamName: streamName,
	}, nil
}

// Publish 发布单个事件
func (p *NATSPublisher) Publish(ctx context.Context, event shared.Event) error {
	body, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	// 使用事件名称作为 Subject
	// 例如: game.player.123.level_up
	subject := event.EventName()

	// 发布到 JetStream（持久化）
	_, err = p.js.Publish(subject, body, nats.Context(ctx))
	if err != nil {
		return fmt.Errorf("failed to publish event to NATS: %w", err)
	}

	log.Printf("Event published to NATS: %s", subject)
	return nil
}

// PublishAll 批量发布事件
func (p *NATSPublisher) PublishAll(ctx context.Context, events []shared.Event) error {
	// NATS 不支持真正的批量发布，但可以流水线发送
	for _, event := range events {
		if err := p.Publish(ctx, event); err != nil {
			return err
		}
	}
	return nil
}

// Close 关闭连接
func (p *NATSPublisher) Close() error {
	if p.nc != nil {
		p.nc.Close()
	}
	return nil
}

// --- 备用方案1：基于 Redis Stream 的实现 ---

// RedisStreamPublisher Redis Stream 事件发布器
// 轻量级备用方案，适合单机或小规模部署
type RedisStreamPublisher struct {
	client    *redis.Client
	streamKey string
}

// NewRedisStreamPublisher 创建 Redis Stream 发布器
// 依赖: go get github.com/redis/go-redis/v9
func NewRedisStreamPublisher(cfg Config) (*RedisStreamPublisher, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	streamKey := cfg.StreamKey
	if streamKey == "" {
		streamKey = "game:events"
	}

	log.Printf("Redis Stream connected successfully (stream: %s)", streamKey)

	return &RedisStreamPublisher{
		client:    client,
		streamKey: streamKey,
	}, nil
}

// Publish 发布单个事件到 Redis Stream
func (p *RedisStreamPublisher) Publish(ctx context.Context, event shared.Event) error {
	body, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	// 使用 XADD 添加事件到 Stream
	// Stream 数据结构: stream_key -> [{id, {event_name, event_data, timestamp}}]
	args := &redis.XAddArgs{
		Stream: p.streamKey,
		MaxLen: 10000, // 限制 Stream 长度，防止无限增长
		Approx: true,  // 近似裁剪，性能更好
		ID:     "*",   // 自动生成 ID
		Values: map[string]interface{}{
			"event_name": event.EventName(),
			"event_data": string(body),
			"timestamp":  time.Now().Unix(),
		},
	}

	id, err := p.client.XAdd(ctx, args).Result()
	if err != nil {
		return fmt.Errorf("failed to publish event to Redis Stream: %w", err)
	}

	log.Printf("Event published to Redis Stream: %s (id: %s)", event.EventName(), id)
	return nil
}

// PublishAll 批量发布事件
func (p *RedisStreamPublisher) PublishAll(ctx context.Context, events []shared.Event) error {
	// Redis Pipeline 批量发送
	pipe := p.client.Pipeline()

	for _, event := range events {
		body, err := json.Marshal(event)
		if err != nil {
			return fmt.Errorf("failed to marshal event: %w", err)
		}

		pipe.XAdd(ctx, &redis.XAddArgs{
			Stream: p.streamKey,
			MaxLen: 10000,
			Approx: true,
			ID:     "*",
			Values: map[string]interface{}{
				"event_name": event.EventName(),
				"event_data": string(body),
				"timestamp":  time.Now().Unix(),
			},
		})
	}

	_, err := pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to publish events to Redis Stream: %w", err)
	}

	log.Printf("Batch published %d events to Redis Stream", len(events))
	return nil
}

// Close 关闭连接
func (p *RedisStreamPublisher) Close() error {
	if p.client != nil {
		return p.client.Close()
	}
	return nil
}

// --- 备用方案2：空实现（开发/测试） ---

// NoopPublisher 空事件发布器（不实际发布，仅打印日志）
// 适用于本地开发、单元测试等场景
type NoopPublisher struct{}

// NewNoopPublisher 创建空发布器
func NewNoopPublisher() *NoopPublisher {
	log.Println("Using NoopPublisher (events will not be actually published)")
	return &NoopPublisher{}
}

// Publish 发布事件（空实现）
func (p *NoopPublisher) Publish(ctx context.Context, event shared.Event) error {
	log.Printf("[NOOP] Event would be published: %s", event.EventName())
	return nil
}

// PublishAll 批量发布事件（空实现）
func (p *NoopPublisher) PublishAll(ctx context.Context, events []shared.Event) error {
	for _, event := range events {
		log.Printf("[NOOP] Event would be published: %s", event.EventName())
	}
	return nil
}

// Close 空实现
func (p *NoopPublisher) Close() error {
	return nil
}
