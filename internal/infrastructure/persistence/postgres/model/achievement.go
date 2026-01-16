// Package model GORM 模型定义
package model

import "time"

// AchievementDefinition 成就定义表
type AchievementDefinition struct {
	BaseModel
	Name           string `gorm:"type:varchar(64);not null;comment:成就名称"`
	Description    string `gorm:"type:text;comment:成就描述"`
	Category       string `gorm:"type:varchar(32);comment:成就分类(pet/social/item/login)"` // pet, social, item, login
	ConditionType  string `gorm:"column:condition_type;type:varchar(32);comment:条件类型"`
	ConditionValue int    `gorm:"column:condition_value;comment:条件数值"`
	RewardCoins    int    `gorm:"column:reward_coins;default:0;comment:奖励金币"`
	RewardDiamonds int    `gorm:"column:reward_diamonds;default:0;comment:奖励钻石"`
	Icon           string `gorm:"type:varchar(64);comment:图标"`
}

// TableName 表名
func (AchievementDefinition) TableName() string {
	return "achievement_definitions"
}

// UserAchievement 用户成就表
type UserAchievement struct {
	BaseModel
	UserID        int       `gorm:"column:user_id;index;not null;comment:用户ID"`
	AchievementID int       `gorm:"column:achievement_id;not null;comment:成就ID"`
	UnlockedAt    time.Time `gorm:"column:unlocked_at;autoCreateTime;comment:解锁时间"`
}

// TableName 表名
func (UserAchievement) TableName() string {
	return "user_achievements"
}
