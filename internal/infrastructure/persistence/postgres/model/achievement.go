// Package model GORM 模型定义
package model

import "time"

// AchievementDefinition 成就定义表
type AchievementDefinition struct {
	BaseModel
	Name           string `gorm:"type:varchar(64);not null"`
	Description    string `gorm:"type:text"`
	Category       string `gorm:"type:varchar(32)"` // pet, social, item, login
	ConditionType  string `gorm:"column:condition_type;type:varchar(32)"`
	ConditionValue int    `gorm:"column:condition_value"`
	RewardCoins    int    `gorm:"column:reward_coins;default:0"`
	RewardDiamonds int    `gorm:"column:reward_diamonds;default:0"`
	Icon           string `gorm:"type:varchar(64)"`
}

// TableName 表名
func (AchievementDefinition) TableName() string {
	return "achievement_definitions"
}

// UserAchievement 用户成就表
type UserAchievement struct {
	BaseModel
	UserID        int       `gorm:"column:user_id;index;not null"`
	AchievementID int       `gorm:"column:achievement_id;not null"`
	UnlockedAt    time.Time `gorm:"column:unlocked_at;autoCreateTime"`
}

// TableName 表名
func (UserAchievement) TableName() string {
	return "user_achievements"
}
