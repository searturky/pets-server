// Package redis 排行榜服务实现
package redis

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"

	"pets-server/internal/application/ranking"
)

// RankingStore 排行榜存储实现
type RankingStore struct {
	client *redis.Client
}

// NewRankingStore 创建排行榜存储
func NewRankingStore(client *redis.Client) *RankingStore {
	return &RankingStore{client: client}
}

func rankingKey(rankType string) string {
	return fmt.Sprintf("ranking:%s", rankType)
}

// GetRanking 获取排行榜
func (r *RankingStore) GetRanking(ctx context.Context, rankType string, offset, limit int) ([]ranking.RankEntry, error) {
	key := rankingKey(rankType)

	// 使用 ZREVRANGE 获取分数从高到低排序的成员
	results, err := r.client.ZRevRangeWithScores(ctx, key, int64(offset), int64(offset+limit-1)).Result()
	if err != nil {
		return nil, err
	}

	entries := make([]ranking.RankEntry, len(results))
	for i, z := range results {
		// Member 是用户ID（存储时需要转为字符串）
		var userID int
		fmt.Sscanf(z.Member.(string), "%d", &userID)
		entries[i] = ranking.RankEntry{
			UserID: userID,
			Score:  int(z.Score),
		}
	}

	return entries, nil
}

// GetUserRank 获取用户排名
func (r *RankingStore) GetUserRank(ctx context.Context, rankType string, userID int) (rank int, score int, err error) {
	key := rankingKey(rankType)
	member := fmt.Sprintf("%d", userID)

	// 获取排名（从0开始，需要+1）
	rankResult, err := r.client.ZRevRank(ctx, key, member).Result()
	if err != nil {
		if err == redis.Nil {
			return 0, 0, nil // 不在排行榜中
		}
		return 0, 0, err
	}

	// 获取分数
	scoreResult, err := r.client.ZScore(ctx, key, member).Result()
	if err != nil {
		return 0, 0, err
	}

	return int(rankResult) + 1, int(scoreResult), nil
}

// UpdateScore 更新用户分数
func (r *RankingStore) UpdateScore(ctx context.Context, rankType string, userID int, score int) error {
	key := rankingKey(rankType)
	member := fmt.Sprintf("%d", userID)

	// 使用 ZADD 更新分数
	return r.client.ZAdd(ctx, key, redis.Z{
		Score:  float64(score),
		Member: member,
	}).Err()
}

// RemoveFromRanking 从排行榜移除用户
func (r *RankingStore) RemoveFromRanking(ctx context.Context, rankType string, userID int) error {
	key := rankingKey(rankType)
	member := fmt.Sprintf("%d", userID)

	return r.client.ZRem(ctx, key, member).Err()
}

// GetRankingCount 获取排行榜总人数
func (r *RankingStore) GetRankingCount(ctx context.Context, rankType string) (int, error) {
	key := rankingKey(rankType)
	count, err := r.client.ZCard(ctx, key).Result()
	return int(count), err
}
