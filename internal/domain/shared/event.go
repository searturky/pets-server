// Package shared 包含跨领域的共享接口和类型
// EventPublisher (事件发布器) 接口定义
// 用于发布领域事件到消息队列
package shared

import "context"

// Event 领域事件基础接口
type Event interface {
	// EventName 返回事件名称，用于路由
	EventName() string
}

// EventPublisher 事件发布器接口
// 定义在领域层，由基础设施层(MQ)实现
type EventPublisher interface {
	// Publish 发布单个事件
	Publish(ctx context.Context, event Event) error

	// PublishAll 批量发布事件
	PublishAll(ctx context.Context, events []Event) error
}

