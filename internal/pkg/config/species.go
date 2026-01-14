// Package config 配置管理
package config

import (
	"github.com/spf13/viper"
)

// SpeciesConfig 物种配置
type SpeciesConfig struct {
	Species []SpeciesEntry `mapstructure:"species"`
	Fusions []FusionEntry  `mapstructure:"fusions"`
}

// SpeciesEntry 物种配置条目
type SpeciesEntry struct {
	ID              int             `mapstructure:"id"`
	Name            string          `mapstructure:"name"`
	Category        string          `mapstructure:"category"` // mammal, avian, fish, reptile, fantasy, elemental
	Rarity          int             `mapstructure:"rarity"`
	IsHidden        bool            `mapstructure:"is_hidden"`
	InterpreterType string          `mapstructure:"interpreter_type"` // feline, canine, slime, etc.
	BaseParts       []string        `mapstructure:"base_parts"`
	SpecialParts    []string        `mapstructure:"special_parts"`
	GenderRule      GenderRuleCfg   `mapstructure:"gender_rule"`
	BreedRules      BreedRulesCfg   `mapstructure:"breed_rules"`
}

// GenderRuleCfg 性别规则配置
type GenderRuleCfg struct {
	Type               string `mapstructure:"type"` // default, asexual, hermaphrodite, mixed
	MaleRatio          int    `mapstructure:"male_ratio"`
	FemaleRatio        int    `mapstructure:"female_ratio"`
	HermaphroditeRatio int    `mapstructure:"hermaphrodite_ratio"`
	CanSelfBreed       bool   `mapstructure:"can_self_breed"`
}

// BreedRulesCfg 繁衍规则配置
type BreedRulesCfg struct {
	MinStage               string `mapstructure:"min_stage"` // egg, child, teen, adult, elderly
	MinLevel               int    `mapstructure:"min_level"`
	MinHappiness           int    `mapstructure:"min_happiness"`
	CooldownHours          int    `mapstructure:"cooldown_hours"`
	SelfBreedCooldownHours int    `mapstructure:"self_breed_cooldown_hours"`
}

// FusionEntry 物种融合配置条目
type FusionEntry struct {
	SpeciesA         int `mapstructure:"species_a"`
	SpeciesB         int `mapstructure:"species_b"`
	Result           int `mapstructure:"result"`
	TriggerThreshold int `mapstructure:"trigger_threshold"`
	Rarity           int `mapstructure:"rarity"`
}

// LoadSpecies 加载物种配置
func LoadSpecies(configPath string) (*SpeciesConfig, error) {
	v := viper.New()
	v.SetConfigFile(configPath)

	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg SpeciesConfig
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
