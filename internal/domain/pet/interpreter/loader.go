// Package interpreter 物种基因解释器
package interpreter

import (
	"fmt"
	"strings"

	"pets-server/internal/domain/pet"
	"pets-server/internal/pkg/config"
)

// BuildSpeciesRegistry 从配置构建物种注册表
func BuildSpeciesRegistry(cfg *config.SpeciesConfig, factory *InterpreterFactory) (*pet.SpeciesRegistry, error) {
	registry := pet.NewSpeciesRegistry()

	for _, entry := range cfg.Species {
		species, err := buildSpecies(entry, factory)
		if err != nil {
			return nil, fmt.Errorf("build species %d (%s): %w", entry.ID, entry.Name, err)
		}
		registry.Register(species)
	}

	return registry, nil
}

// BuildFusionRegistry 从配置构建融合注册表
func BuildFusionRegistry(cfg *config.SpeciesConfig) *pet.SpeciesFusionRegistry {
	registry := pet.NewSpeciesFusionRegistry()

	for _, fusion := range cfg.Fusions {
		registry.Register(
			pet.SpeciesID(fusion.SpeciesA),
			pet.SpeciesID(fusion.SpeciesB),
			pet.HiddenSpeciesConfig{
				ResultSpecies:    pet.SpeciesID(fusion.Result),
				TriggerThreshold: fusion.TriggerThreshold,
				Rarity:           fusion.Rarity,
			},
		)
	}

	return registry
}

// buildSpecies 从配置条目构建物种对象
func buildSpecies(entry config.SpeciesEntry, factory *InterpreterFactory) (*pet.Species, error) {
	// 获取解释器
	var interpreter pet.GeneInterpreter
	if entry.InterpreterType != "" {
		interp, ok := factory.Get(entry.InterpreterType)
		if !ok {
			return nil, fmt.Errorf("unknown interpreter type: %s", entry.InterpreterType)
		}
		interpreter = interp
	}

	// 解析分类
	category := parseCategory(entry.Category)

	// 解析部位
	baseParts := parseParts(entry.BaseParts)
	specialParts := parseParts(entry.SpecialParts)

	// 解析性别规则
	genderRule := parseGenderRule(entry.GenderRule)

	// 解析繁衍规则
	breedRules := parseBreedRules(entry.BreedRules)

	return &pet.Species{
		ID:           pet.SpeciesID(entry.ID),
		Name:         entry.Name,
		Category:     category,
		BaseParts:    baseParts,
		SpecialParts: specialParts,
		Rarity:       entry.Rarity,
		IsHidden:     entry.IsHidden,
		GenderRule:   genderRule,
		BreedRules:   breedRules,
		Interpreter:  interpreter,
	}, nil
}

// parseCategory 解析物种分类
func parseCategory(s string) pet.SpeciesCategory {
	switch strings.ToLower(s) {
	case "mammal":
		return pet.CategoryMammal
	case "avian":
		return pet.CategoryAvian
	case "fish":
		return pet.CategoryFish
	case "reptile":
		return pet.CategoryReptile
	case "fantasy":
		return pet.CategoryFantasy
	case "elemental":
		return pet.CategoryElemental
	default:
		return pet.CategoryUnknown
	}
}

// parseParts 解析部位列表
func parseParts(parts []string) []pet.PartType {
	result := make([]pet.PartType, 0, len(parts))
	for _, p := range parts {
		partType := parsePartType(p)
		result = append(result, partType)
	}
	return result
}

// parsePartType 解析单个部位类型
func parsePartType(s string) pet.PartType {
	switch strings.ToLower(s) {
	case "none":
		return pet.PartTypeNone
	case "ear":
		return pet.PartTypeEar
	case "tail":
		return pet.PartTypeTail
	case "fur":
		return pet.PartTypeFur
	case "wing":
		return pet.PartTypeWing
	case "beak":
		return pet.PartTypeBeak
	case "crest":
		return pet.PartTypeCrest
	case "fin":
		return pet.PartTypeFin
	case "scale":
		return pet.PartTypeScale
	case "tail_fin":
		return pet.PartTypeTailFin
	case "shell":
		return pet.PartTypeShell
	case "horn":
		return pet.PartTypeHorn
	case "armor":
		return pet.PartTypeArmor
	case "aura":
		return pet.PartTypeAura
	case "claw":
		return pet.PartTypeClaw
	case "whisker":
		return pet.PartTypeWhisker
	default:
		return pet.PartTypeNone
	}
}

// parseGenderRule 解析性别规则
func parseGenderRule(cfg config.GenderRuleCfg) pet.GenderRule {
	switch strings.ToLower(cfg.Type) {
	case "asexual":
		rule := pet.AsexualGenderRule()
		rule.CanSelfBreed = cfg.CanSelfBreed
		return rule
	case "hermaphrodite":
		return pet.HermaphroditeGenderRule()
	case "mixed":
		return pet.MixedGenderRule(cfg.MaleRatio, cfg.FemaleRatio, cfg.HermaphroditeRatio)
	default: // "default" or empty
		return pet.DefaultGenderRule()
	}
}

// parseBreedRules 解析繁衍规则
func parseBreedRules(cfg config.BreedRulesCfg) pet.BreedingRules {
	// 使用默认值
	rules := pet.DefaultBreedingRules()

	// 覆盖配置中的值
	if cfg.MinStage != "" {
		rules.MinStage = parseStage(cfg.MinStage)
	}
	if cfg.MinLevel > 0 {
		rules.MinLevel = cfg.MinLevel
	}
	if cfg.MinHappiness > 0 {
		rules.MinHappiness = cfg.MinHappiness
	}
	if cfg.CooldownHours > 0 {
		rules.CooldownHours = cfg.CooldownHours
	}
	if cfg.SelfBreedCooldownHours > 0 {
		rules.SelfBreedCooldownHours = cfg.SelfBreedCooldownHours
	}

	return rules
}

// parseStage 解析成长阶段
func parseStage(s string) pet.Stage {
	switch strings.ToLower(s) {
	case "egg":
		return pet.StageEgg
	case "child":
		return pet.StageChild
	case "teen":
		return pet.StageTeen
	case "adult":
		return pet.StageAdult
	case "elderly":
		return pet.StageElderly
	default:
		return pet.StageAdult
	}
}
