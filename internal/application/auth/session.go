// Package auth 认证应用服务
package auth

import (
	"context"
	"time"
)

// SessionStore 会话存储接口（在应用层定义，基础设施层实现）
// 用于实现单设备登录的会话校验。
type SessionStore interface {
	// SetCurrentSession 设置当前有效会话
	SetCurrentSession(ctx context.Context, userID int, sessionID string, ttl time.Duration) error
	// GetCurrentSession 获取当前有效会话
	GetCurrentSession(ctx context.Context, userID int) (string, error)
	// DeleteCurrentSession 删除当前有效会话
	DeleteCurrentSession(ctx context.Context, userID int) error
}
