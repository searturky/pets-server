// Package redis 会话存储实现
package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// AuthSessionStore 认证会话存储
type AuthSessionStore struct {
	client *redis.Client
}

// NewAuthSessionStore 创建认证会话存储
func NewAuthSessionStore(client *redis.Client) *AuthSessionStore {
	return &AuthSessionStore{client: client}
}

func authSessionKey(userID int) string {
	return fmt.Sprintf("pets:auth:session:%d", userID)
}

// SetCurrentSession 设置当前有效会话
func (s *AuthSessionStore) SetCurrentSession(ctx context.Context, userID int, sessionID string, ttl time.Duration) error {
	return s.client.Set(ctx, authSessionKey(userID), sessionID, ttl).Err()
}

// GetCurrentSession 获取当前有效会话
func (s *AuthSessionStore) GetCurrentSession(ctx context.Context, userID int) (string, error) {
	val, err := s.client.Get(ctx, authSessionKey(userID)).Result()
	if err == redis.Nil {
		return "", nil
	}
	return val, err
}

// DeleteCurrentSession 删除当前有效会话
func (s *AuthSessionStore) DeleteCurrentSession(ctx context.Context, userID int) error {
	return s.client.Del(ctx, authSessionKey(userID)).Err()
}
