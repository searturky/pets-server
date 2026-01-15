// Package ranking 排行榜应用服务
// 处理排行榜相关的业务用例
package ranking

import (
	"context"
)

// RankingStore 排行榜存储接口（由 Redis 实现）
type RankingStore interface {
	// GetRanking 获取排行榜
	GetRanking(ctx context.Context, rankType string, offset, limit int) ([]RankEntry, error)

	// GetUserRank 获取用户排名
	GetUserRank(ctx context.Context, rankType string, userID int) (rank int, score int, err error)

	// UpdateScore 更新用户分数
	UpdateScore(ctx context.Context, rankType string, userID int, score int) error
}

// RankEntry 排行条目
type RankEntry struct {
	UserID int
	Score  int
}

// Service 排行榜应用服务
type Service struct {
	store RankingStore
}

// NewService 创建排行榜服务
func NewService(store RankingStore) *Service {
	return &Service{store: store}
}

// GetRanking 获取排行榜
func (s *Service) GetRanking(ctx context.Context, userID int, req RankingRequest) (*RankingResponse, error) {
	if req.Limit <= 0 {
		req.Limit = 20
	}
	if req.Limit > 100 {
		req.Limit = 100
	}

	// 获取排行榜数据
	entries, err := s.store.GetRanking(ctx, string(req.Type), req.Offset, req.Limit)
	if err != nil {
		return nil, err
	}

	// 转换为 DTO
	rankings := make([]RankItemDTO, 0, len(entries))
	for i, entry := range entries {
		rankings = append(rankings, RankItemDTO{
			Rank:   req.Offset + i + 1,
			UserID: entry.UserID,
			Score:  entry.Score,
			// TODO: 批量获取用户昵称、头像等信息
		})
	}

	// 获取当前用户排名
	var myRank *RankItemDTO
	rank, score, err := s.store.GetUserRank(ctx, string(req.Type), userID)
	if err == nil && rank > 0 {
		myRank = &RankItemDTO{
			Rank:   rank,
			UserID: userID,
			Score:  score,
		}
	}

	return &RankingResponse{
		Type:     req.Type,
		Rankings: rankings,
		MyRank:   myRank,
	}, nil
}

// UpdatePetLevelRank 更新宠物等级排行
func (s *Service) UpdatePetLevelRank(ctx context.Context, userID int, level int) error {
	return s.store.UpdateScore(ctx, string(RankTypePetLevel), userID, level)
}

// UpdateAchievementRank 更新成就排行
func (s *Service) UpdateAchievementRank(ctx context.Context, userID int, count int) error {
	return s.store.UpdateScore(ctx, string(RankTypeAchievement), userID, count)
}

