// Package model GORM 模型定义
package model

import "time"

// Pet 宠物表
type Pet struct {
	BaseModel
	UserID int    `gorm:"index;not null;comment:用户ID"` // 一个用户可以有多只宠物
	Name   string `gorm:"type:varchar(32);not null;comment:宠物名称"`

	// 物种和性别
	SpeciesID int   `gorm:"column:species_id;index;not null;default:101;comment:物种ID"` // 物种ID，默认为猫
	Gender    int16 `gorm:"column:gender;default:1;comment:性别(0无1雄2雌3雌雄同体)"`           // 0=无 1=雄 2=雌 3=雌雄同体

	// 基因系统 (40位十六进制)
	GeneCode string `gorm:"column:gene_code;type:varchar(40);not null;comment:基因编码(40位十六进制)"`

	// 通用外观属性 (从基因解析并缓存)
	ColorPrimary   string `gorm:"column:color_primary;type:varchar(7);comment:主色(十六进制颜色)"`
	ColorSecondary string `gorm:"column:color_secondary;type:varchar(7);comment:副色(十六进制颜色)"`
	PatternType    int16  `gorm:"column:pattern_type;comment:花纹类型"`
	PatternDensity int16  `gorm:"column:pattern_density;comment:花纹密度"`
	BodyType       int16  `gorm:"column:body_type;comment:体型"`
	EyeShape       int16  `gorm:"column:eye_shape;comment:眼睛形状"`
	EyeColor       int16  `gorm:"column:eye_color;comment:眼睛颜色"`

	// 物种特有外观 (JSON存储)
	SpecialAppearance string `gorm:"column:special_appearance;type:jsonb;comment:物种特有外观(JSON)"`

	// 性格属性 (从基因解析, 0-100)
	TraitActivity     int16 `gorm:"column:trait_activity;comment:活跃度(0-100)"`
	TraitAppetite     int16 `gorm:"column:trait_appetite;comment:食欲(0-100)"`
	TraitSocial       int16 `gorm:"column:trait_social;comment:社交性(0-100)"`
	TraitCuriosity    int16 `gorm:"column:trait_curiosity;comment:好奇心(0-100)"`
	TraitTemper       int16 `gorm:"column:trait_temper;comment:脾气(0-100)"`
	TraitLoyalty      int16 `gorm:"column:trait_loyalty;comment:忠诚度(0-100)"`
	TraitIntelligence int16 `gorm:"column:trait_intelligence;comment:智力(0-100)"`
	TraitPlayfulness  int16 `gorm:"column:trait_playfulness;comment:玩耍性(0-100)"`

	// 技能
	SkillID          int   `gorm:"column:skill_id;comment:主技能ID"`
	SkillLevel       int16 `gorm:"column:skill_level;default:1;comment:技能等级"`
	SkillStrength    int16 `gorm:"column:skill_strength;comment:技能强度"`
	SkillSecondaryID int   `gorm:"column:skill_secondary_id;comment:副技能ID"`

	// 成长状态
	Stage int16 `gorm:"default:0;comment:成长阶段(0蛋1幼年2成长3成熟4老年)"` // 0蛋 1幼年 2成长 3成熟 4老年
	Exp   int   `gorm:"default:0;comment:经验值"`
	Level int   `gorm:"default:1;comment:等级"`

	// 实时状态 (0-100)
	Hunger      int16 `gorm:"default:50;comment:饥饿度(0-100)"`
	Happiness   int16 `gorm:"default:50;comment:快乐度(0-100)"`
	Cleanliness int16 `gorm:"default:50;comment:清洁度(0-100)"`
	Energy      int16 `gorm:"default:100;comment:能量值(0-100)"`

	// 繁衍相关
	Parent1ID   *int       `gorm:"column:parent1_id;index;comment:父方ID"` // 父方ID (可为空)
	Parent2ID   *int       `gorm:"column:parent2_id;index;comment:母方ID"` // 母方ID (可为空)
	Generation  int        `gorm:"column:generation;default:0;comment:代数"`
	LastBreedAt *time.Time `gorm:"column:last_breed_at;comment:最后繁殖时间"`

	// 时间记录
	LastFedAt     time.Time `gorm:"column:last_fed_at;comment:最后喂食时间"`
	LastPlayedAt  time.Time `gorm:"column:last_played_at;comment:最后玩耍时间"`
	LastCleanedAt time.Time `gorm:"column:last_cleaned_at;comment:最后清洁时间"`
	BornAt        time.Time `gorm:"column:born_at;comment:出生时间"`
}

// TableName 表名
func (Pet) TableName() string {
	return "pets"
}

// SpeciesDefinition 物种定义表
type SpeciesDefinition struct {
	BaseModel
	Name         string `gorm:"type:varchar(32);not null;comment:物种名称"`
	Category     string `gorm:"type:varchar(16);not null;comment:物种分类(mammal/avian/fish/reptile/fantasy/elemental)"` // mammal/avian/fish/reptile/fantasy/elemental
	BaseParts    string `gorm:"type:jsonb;not null;comment:基础部位(JSON数组)"`                                            // JSON: ["body", "eye", "pattern"]
	SpecialParts string `gorm:"type:jsonb;not null;comment:特殊部位(JSON数组)"`                                            // JSON: ["ear", "tail"] 或 ["wing", "beak"]
	GeneMapping  string `gorm:"type:jsonb;not null;comment:基因映射规则(JSON)"`                                            // JSON: 基因位置到特征的映射
	Rarity       int    `gorm:"default:1;comment:稀有度(1-5)"`                                                          // 稀有度 1-5
	GenderRules  string `gorm:"type:jsonb;not null;comment:性别规则(JSON)"`                                              // JSON: {"allowed": [1,2], "ratio": {"1":50,"2":50}, "can_self_breed": false}
	BreedRules   string `gorm:"type:jsonb;comment:繁殖规则(JSON)"`                                                       // JSON: 繁殖规则
	IsHidden     bool   `gorm:"default:false;comment:是否隐藏物种"`                                                        // 是否为隐藏物种
}

// TableName 表名
func (SpeciesDefinition) TableName() string {
	return "species_definitions"
}

// SpeciesFusion 物种融合表
type SpeciesFusion struct {
	BaseModel
	SpeciesAID       int `gorm:"column:species_a_id;not null;index;comment:物种A的ID"`
	SpeciesBID       int `gorm:"column:species_b_id;not null;index;comment:物种B的ID"`
	ResultSpeciesID  int `gorm:"column:result_species_id;not null;comment:融合结果物种ID"`
	TriggerThreshold int `gorm:"column:trigger_threshold;default:200;comment:触发阈值"`
}

// TableName 表名
func (SpeciesFusion) TableName() string {
	return "species_fusions"
}

// SkillDefinition 技能定义表
type SkillDefinition struct {
	BaseModel
	Name        string `gorm:"type:varchar(32);not null;comment:技能名称"`
	Description string `gorm:"type:text;comment:技能描述"`
	EffectType  string `gorm:"column:effect_type;type:varchar(32);comment:效果类型"`
	EffectValue int    `gorm:"column:effect_value;comment:效果数值"`
	Rarity      int16  `gorm:"default:1;comment:稀有度(1普通2稀有3史诗4传说)"` // 1普通 2稀有 3史诗 4传说
}

// TableName 表名
func (SkillDefinition) TableName() string {
	return "skill_definitions"
}
