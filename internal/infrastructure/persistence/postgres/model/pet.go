// Package model GORM 模型定义
package model

import "time"

// Pet 宠物表
type Pet struct {
	BaseModel
	UserID int `gorm:"index;not null"` // 一个用户可以有多只宠物
	Name   string `gorm:"type:varchar(32);not null"`

	// 物种和性别
	SpeciesID int   `gorm:"column:species_id;index;not null;default:101"` // 物种ID，默认为猫
	Gender    int16 `gorm:"column:gender;default:1"`                      // 0=无 1=雄 2=雌 3=雌雄同体

	// 基因系统 (40位十六进制)
	GeneCode string `gorm:"column:gene_code;type:varchar(40);not null"`

	// 通用外观属性 (从基因解析并缓存)
	ColorPrimary   string `gorm:"column:color_primary;type:varchar(7)"`
	ColorSecondary string `gorm:"column:color_secondary;type:varchar(7)"`
	PatternType    int16  `gorm:"column:pattern_type"`
	PatternDensity int16  `gorm:"column:pattern_density"`
	BodyType       int16  `gorm:"column:body_type"`
	EyeShape       int16  `gorm:"column:eye_shape"`
	EyeColor       int16  `gorm:"column:eye_color"`

	// 物种特有外观 (JSON存储)
	SpecialAppearance string `gorm:"column:special_appearance;type:jsonb"`

	// 性格属性 (从基因解析, 0-100)
	TraitActivity     int16 `gorm:"column:trait_activity"`
	TraitAppetite     int16 `gorm:"column:trait_appetite"`
	TraitSocial       int16 `gorm:"column:trait_social"`
	TraitCuriosity    int16 `gorm:"column:trait_curiosity"`
	TraitTemper       int16 `gorm:"column:trait_temper"`
	TraitLoyalty      int16 `gorm:"column:trait_loyalty"`
	TraitIntelligence int16 `gorm:"column:trait_intelligence"`
	TraitPlayfulness  int16 `gorm:"column:trait_playfulness"`

	// 技能
	SkillID          int   `gorm:"column:skill_id"`
	SkillLevel       int16 `gorm:"column:skill_level;default:1"`
	SkillStrength    int16 `gorm:"column:skill_strength"`
	SkillSecondaryID int   `gorm:"column:skill_secondary_id"`

	// 成长状态
	Stage int16 `gorm:"default:0"` // 0蛋 1幼年 2成长 3成熟 4老年
	Exp   int   `gorm:"default:0"`
	Level int   `gorm:"default:1"`

	// 实时状态 (0-100)
	Hunger      int16 `gorm:"default:50"`
	Happiness   int16 `gorm:"default:50"`
	Cleanliness int16 `gorm:"default:50"`
	Energy      int16 `gorm:"default:100"`

	// 繁衍相关
	Parent1ID   *int       `gorm:"column:parent1_id;index"` // 父方ID (可为空)
	Parent2ID   *int       `gorm:"column:parent2_id;index"` // 母方ID (可为空)
	Generation  int        `gorm:"column:generation;default:0"`
	LastBreedAt *time.Time `gorm:"column:last_breed_at"`

	// 时间记录
	LastFedAt     time.Time `gorm:"column:last_fed_at"`
	LastPlayedAt  time.Time `gorm:"column:last_played_at"`
	LastCleanedAt time.Time `gorm:"column:last_cleaned_at"`
	BornAt        time.Time `gorm:"column:born_at"`
}

// TableName 表名
func (Pet) TableName() string {
	return "pets"
}

// SpeciesDefinition 物种定义表
type SpeciesDefinition struct {
	BaseModel
	Name         string `gorm:"type:varchar(32);not null"`
	Category     string `gorm:"type:varchar(16);not null"` // mammal/avian/fish/reptile/fantasy/elemental
	BaseParts    string `gorm:"type:jsonb;not null"`       // JSON: ["body", "eye", "pattern"]
	SpecialParts string `gorm:"type:jsonb;not null"`       // JSON: ["ear", "tail"] 或 ["wing", "beak"]
	GeneMapping  string `gorm:"type:jsonb;not null"`       // JSON: 基因位置到特征的映射
	Rarity       int    `gorm:"default:1"`                 // 稀有度 1-5
	GenderRules  string `gorm:"type:jsonb;not null"`       // JSON: {"allowed": [1,2], "ratio": {"1":50,"2":50}, "can_self_breed": false}
	BreedRules   string `gorm:"type:jsonb"`                // JSON: 繁殖规则
	IsHidden     bool   `gorm:"default:false"`             // 是否为隐藏物种
}

// TableName 表名
func (SpeciesDefinition) TableName() string {
	return "species_definitions"
}

// SpeciesFusion 物种融合表
type SpeciesFusion struct {
	BaseModel
	SpeciesAID       int `gorm:"column:species_a_id;not null;index"`
	SpeciesBID       int `gorm:"column:species_b_id;not null;index"`
	ResultSpeciesID  int `gorm:"column:result_species_id;not null"`
	TriggerThreshold int `gorm:"column:trigger_threshold;default:200"`
}

// TableName 表名
func (SpeciesFusion) TableName() string {
	return "species_fusions"
}

// SkillDefinition 技能定义表
type SkillDefinition struct {
	BaseModel
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
