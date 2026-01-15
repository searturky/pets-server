// Package achievement 成就领域
// Repository 仓储接口
package achievement

import "context"

// Repository 成就仓储接口
type Repository interface {
	// --- 用户成就 ---
	
	// FindByUserID 获取用户所有成就
	FindByUserID(ctx context.Context, userID int) ([]*UserAchievement, error)

	// HasAchievement 检查用户是否已获得某成就
	HasAchievement(ctx context.Context, userID int, achievementID int) (bool, error)

	// Save 保存用户成就
	Save(ctx context.Context, achievement *UserAchievement) error

	// --- 成就定义 ---
	
	// GetDefinition 获取成就定义
	GetDefinition(ctx context.Context, id int) (*AchievementDefinition, error)

	// GetAllDefinitions 获取所有成就定义
	GetAllDefinitions(ctx context.Context) ([]*AchievementDefinition, error)

	// GetByCategory 获取某分类的所有成就定义
	GetByCategory(ctx context.Context, category AchievementCategory) ([]*AchievementDefinition, error)
}

