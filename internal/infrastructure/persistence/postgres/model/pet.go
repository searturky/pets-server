// Package model GORM 模型定义
package model

import "time"

// Pet 宠物表
type Pet struct {
	ID     int64  `gorm:"primaryKey;autoIncrement"`
	UserID int64  `gorm:"uniqueIndex;not null"` // 一个用户只能有一只宠物
	Name   string `gorm:"type:varchar(32);not null"`

	// 基因系统 (20位十六进制)
	GeneCode string `gorm:"column:gene_code;type:varchar(20);not null"`

	// 外观属性 (从基因解析并缓存)
	ColorPrimary   string `gorm:"column:color_primary;type:varchar(7)"`
	ColorSecondary string `gorm:"column:color_secondary;type:varchar(7)"`
	PatternType    int16  `gorm:"column:pattern_type"`
	BodyType       int16  `gorm:"column:body_type"`
	EarType        int16  `gorm:"column:ear_type"`
	TailType       int16  `gorm:"column:tail_type"`
	EyeType        int16  `gorm:"column:eye_type"`

	// 性格属性 (从基因解析, 0-100)
	TraitActivity  int16 `gorm:"column:trait_activity"`
	TraitAppetite  int16 `gorm:"column:trait_appetite"`
	TraitSocial    int16 `gorm:"column:trait_social"`
	TraitCuriosity int16 `gorm:"column:trait_curiosity"`

	// 技能
	SkillID       int   `gorm:"column:skill_id"`
	SkillLevel    int16 `gorm:"column:skill_level;default:1"`
	SkillStrength int16 `gorm:"column:skill_strength"`

	// 成长状态
	Stage int16 `gorm:"default:0"` // 0蛋 1幼年 2成长 3成熟 4老年
	Exp   int   `gorm:"default:0"`
	Level int   `gorm:"default:1"`

	// 实时状态 (0-100)
	Hunger      int16 `gorm:"default:50"`
	Happiness   int16 `gorm:"default:50"`
	Cleanliness int16 `gorm:"default:50"`
	Energy      int16 `gorm:"default:100"`

	// 时间记录
	LastFedAt     time.Time `gorm:"column:last_fed_at"`
	LastPlayedAt  time.Time `gorm:"column:last_played_at"`
	LastCleanedAt time.Time `gorm:"column:last_cleaned_at"`
	BornAt        time.Time `gorm:"column:born_at"`
	CreatedAt     time.Time `gorm:"autoCreateTime"`
}

// TableName 表名
func (Pet) TableName() string {
	return "pets"
}

// SkillDefinition 技能定义表
type SkillDefinition struct {
	ID          int    `gorm:"primaryKey"`
	Name        string `gorm:"type:varchar(32);not null"`
	Description string `gorm:"type:text"`
	EffectType  string `gorm:"column:effect_type;type:varchar(32)"`
	EffectValue int    `gorm:"column:effect_value"`
	Rarity      int16  `gorm:"default:1"` // 1普通 2稀有 3史诗 4传说
}

// TableName 表名
func (SkillDefinition) TableName() string {
	return "skill_definitions"
}

