// Package redis 缓存服务实现
package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	petApp "pets-server/internal/application/pet"
)

// CacheService 缓存服务实现
type CacheService struct {
	client *redis.Client
}

// NewCacheService 创建缓存服务
func NewCacheService(client *redis.Client) *CacheService {
	return &CacheService{client: client}
}

// --- 宠物详情缓存 ---

func petDetailKey(userID int64) string {
	return fmt.Sprintf("pet:detail:%d", userID)
}

// GetPetDetail 获取宠物详情缓存
func (c *CacheService) GetPetDetail(ctx context.Context, userID int64) (*petApp.PetDetailDTO, error) {
	key := petDetailKey(userID)

	data, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // 缓存未命中
		}
		return nil, err
	}

	var dto petApp.PetDetailDTO
	if err := json.Unmarshal(data, &dto); err != nil {
		return nil, err
	}

	return &dto, nil
}

// SetPetDetail 设置宠物详情缓存
func (c *CacheService) SetPetDetail(ctx context.Context, userID int64, pet *petApp.PetDetailDTO, ttl time.Duration) error {
	key := petDetailKey(userID)

	data, err := json.Marshal(pet)
	if err != nil {
		return err
	}

	return c.client.Set(ctx, key, data, ttl).Err()
}

// DeletePetDetail 删除宠物详情缓存
func (c *CacheService) DeletePetDetail(ctx context.Context, userID int64) error {
	key := petDetailKey(userID)
	return c.client.Del(ctx, key).Err()
}

// --- 通用缓存方法 ---

// Get 获取缓存
func (c *CacheService) Get(ctx context.Context, key string) (string, error) {
	return c.client.Get(ctx, key).Result()
}

// Set 设置缓存
func (c *CacheService) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	return c.client.Set(ctx, key, value, ttl).Err()
}

// Delete 删除缓存
func (c *CacheService) Delete(ctx context.Context, key string) error {
	return c.client.Del(ctx, key).Err()
}

// Exists 检查键是否存在
func (c *CacheService) Exists(ctx context.Context, key string) (bool, error) {
	n, err := c.client.Exists(ctx, key).Result()
	return n > 0, err
}
