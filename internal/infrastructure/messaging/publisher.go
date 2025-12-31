// Package messaging 消息队列事件发布实现
package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"

	"pets-server/internal/domain/shared"
)

// Config MQ 配置
type Config struct {
	URL      string // amqp://guest:guest@localhost:5672/
	Exchange string
	Queue    string
}

// RabbitMQPublisher RabbitMQ 事件发布器实现
type RabbitMQPublisher struct {
	conn     *amqp.Connection
	channel  *amqp.Channel
	exchange string
}

// NewRabbitMQPublisher 创建 RabbitMQ 发布器
func NewRabbitMQPublisher(cfg Config) (*RabbitMQPublisher, error) {
	conn, err := amqp.Dial(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	// 声明交换机
	err = ch.ExchangeDeclare(
		cfg.Exchange, // name
		"topic",      // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to declare exchange: %w", err)
	}

	log.Println("RabbitMQ connected successfully")

	return &RabbitMQPublisher{
		conn:     conn,
		channel:  ch,
		exchange: cfg.Exchange,
	}, nil
}

// Publish 发布单个事件
func (p *RabbitMQPublisher) Publish(ctx context.Context, event shared.Event) error {
	body, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	// 使用事件名称作为路由键
	routingKey := event.EventName()

	err = p.channel.PublishWithContext(
		ctx,
		p.exchange, // exchange
		routingKey, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent, // 持久化消息
			Timestamp:    time.Now(),
		},
	)

	if err != nil {
		return fmt.Errorf("failed to publish event: %w", err)
	}

	log.Printf("Event published: %s", routingKey)
	return nil
}

// PublishAll 批量发布事件
func (p *RabbitMQPublisher) PublishAll(ctx context.Context, events []shared.Event) error {
	for _, event := range events {
		if err := p.Publish(ctx, event); err != nil {
			return err
		}
	}
	return nil
}

// Close 关闭连接
func (p *RabbitMQPublisher) Close() error {
	if p.channel != nil {
		p.channel.Close()
	}
	if p.conn != nil {
		p.conn.Close()
	}
	return nil
}

// --- 备用：基于 Redis Stream 的实现 ---

// RedisStreamPublisher Redis Stream 事件发布器
// 如果不想使用 RabbitMQ，可以用 Redis Stream 作为简单替代
type RedisStreamPublisher struct {
	// TODO: 实现基于 Redis Stream 的事件发布
	// 使用 XADD 命令添加事件到 Stream
}

// --- 空实现（用于开发/测试） ---

// NoopPublisher 空事件发布器（不实际发布）
type NoopPublisher struct{}

// NewNoopPublisher 创建空发布器
func NewNoopPublisher() *NoopPublisher {
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

