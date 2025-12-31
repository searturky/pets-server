// Package repo 仓储实现
package repo

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"pets-server/internal/domain/achievement"
	"pets-server/internal/infrastructure/persistence/postgres"
	"pets-server/internal/infrastructure/persistence/postgres/model"
)

// AchievementRepository 成就仓储实现
type AchievementRepository struct {
	db *gorm.DB
}

// NewAchievementRepository 创建成就仓储
func NewAchievementRepository(db *gorm.DB) *AchievementRepository {
	return &AchievementRepository{db: db}
}

// FindByUserID 获取用户所有成就
func (r *AchievementRepository) FindByUserID(ctx context.Context, userID int64) ([]*achievement.UserAchievement, error) {
	db := postgres.GetTx(ctx, r.db)

	var models []model.UserAchievement
	if err := db.Where("user_id = ?", userID).Find(&models).Error; err != nil {
		return nil, err
	}

	achievements := make([]*achievement.UserAchievement, len(models))
	for i, m := range models {
		achievements[i] = &achievement.UserAchievement{
			ID:            m.ID,
			UserID:        m.UserID,
			AchievementID: m.AchievementID,
			UnlockedAt:    m.UnlockedAt,
		}
	}

	return achievements, nil
}

// HasAchievement 检查用户是否已获得某成就
func (r *AchievementRepository) HasAchievement(ctx context.Context, userID int64, achievementID int) (bool, error) {
	db := postgres.GetTx(ctx, r.db)

	var count int64
	if err := db.Model(&model.UserAchievement{}).
		Where("user_id = ? AND achievement_id = ?", userID, achievementID).
		Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

// Save 保存用户成就
func (r *AchievementRepository) Save(ctx context.Context, a *achievement.UserAchievement) error {
	db := postgres.GetTx(ctx, r.db)

	m := &model.UserAchievement{
		ID:            a.ID,
		UserID:        a.UserID,
		AchievementID: a.AchievementID,
		UnlockedAt:    a.UnlockedAt,
	}

	if err := db.Create(m).Error; err != nil {
		return err
	}

	a.ID = m.ID
	return nil
}

// GetDefinition 获取成就定义
func (r *AchievementRepository) GetDefinition(ctx context.Context, id int) (*achievement.AchievementDefinition, error) {
	db := postgres.GetTx(ctx, r.db)

	var m model.AchievementDefinition
	if err := db.First(&m, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return r.definitionToDomain(&m), nil
}

// GetAllDefinitions 获取所有成就定义
func (r *AchievementRepository) GetAllDefinitions(ctx context.Context) ([]*achievement.AchievementDefinition, error) {
	db := postgres.GetTx(ctx, r.db)

	var models []model.AchievementDefinition
	if err := db.Find(&models).Error; err != nil {
		return nil, err
	}

	defs := make([]*achievement.AchievementDefinition, len(models))
	for i, m := range models {
		defs[i] = r.definitionToDomain(&m)
	}

	return defs, nil
}

// GetByCategory 获取某分类的所有成就定义
func (r *AchievementRepository) GetByCategory(ctx context.Context, category achievement.AchievementCategory) ([]*achievement.AchievementDefinition, error) {
	db := postgres.GetTx(ctx, r.db)

	var models []model.AchievementDefinition
	if err := db.Where("category = ?", string(category)).Find(&models).Error; err != nil {
		return nil, err
	}

	defs := make([]*achievement.AchievementDefinition, len(models))
	for i, m := range models {
		defs[i] = r.definitionToDomain(&m)
	}

	return defs, nil
}

func (r *AchievementRepository) definitionToDomain(m *model.AchievementDefinition) *achievement.AchievementDefinition {
	return &achievement.AchievementDefinition{
		ID:             m.ID,
		Name:           m.Name,
		Description:    m.Description,
		Category:       achievement.AchievementCategory(m.Category),
		ConditionType:  achievement.ConditionType(m.ConditionType),
		ConditionValue: m.ConditionValue,
		RewardCoins:    m.RewardCoins,
		RewardDiamonds: m.RewardDiamonds,
		Icon:           m.Icon,
	}
}
