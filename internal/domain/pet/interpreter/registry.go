// Package interpreter 物种基因解释器
// 注册表初始化 - 注册所有物种和解释器
package interpreter

import "pets-server/internal/domain/pet"

// InitDefaultRegistry 初始化默认物种注册表
func InitDefaultRegistry() *pet.SpeciesRegistry {
	registry := pet.NewSpeciesRegistry()

	// 注册哺乳类
	registerMammals(registry)

	// 注册鸟类
	registerAvians(registry)

	// 注册水生类
	registerAquatics(registry)

	// 注册幻想类
	registerFantasy(registry)

	return registry
}

// InitDefaultFusionRegistry 初始化默认物种融合注册表
func InitDefaultFusionRegistry() *pet.SpeciesFusionRegistry {
	registry := pet.NewSpeciesFusionRegistry()

	// === 格里芬相关融合 ===
	// 猫 + 鹦鹉 = 格里芬
	registry.Register(pet.SpeciesCat, pet.SpeciesParrot, pet.HiddenSpeciesConfig{
		ResultSpecies:    pet.SpeciesGriffin,
		TriggerThreshold: 200,
		Rarity:           5,
	})
	// 猫 + 猫头鹰 = 格里芬
	registry.Register(pet.SpeciesCat, pet.SpeciesOwl, pet.HiddenSpeciesConfig{
		ResultSpecies:    pet.SpeciesGriffin,
		TriggerThreshold: 190,
		Rarity:           5,
	})
	// 狗 + 鹦鹉 = 格里芬
	registry.Register(pet.SpeciesDog, pet.SpeciesParrot, pet.HiddenSpeciesConfig{
		ResultSpecies:    pet.SpeciesGriffin,
		TriggerThreshold: 210,
		Rarity:           5,
	})

	// === 龙相关融合 ===
	// 金鱼 + 蜥蜴 = 龙
	registry.Register(pet.SpeciesGoldfish, pet.SpeciesLizard, pet.HiddenSpeciesConfig{
		ResultSpecies:    pet.SpeciesDragon,
		TriggerThreshold: 220,
		Rarity:           5,
	})
	// 热带鱼 + 蜥蜴 = 龙
	registry.Register(pet.SpeciesTropicalFish, pet.SpeciesLizard, pet.HiddenSpeciesConfig{
		ResultSpecies:    pet.SpeciesDragon,
		TriggerThreshold: 210,
		Rarity:           5,
	})
	// 凤凰 + 蜥蜴 = 龙
	registry.Register(pet.SpeciesPhoenix, pet.SpeciesLizard, pet.HiddenSpeciesConfig{
		ResultSpecies:    pet.SpeciesDragon,
		TriggerThreshold: 180,
		Rarity:           5,
	})

	// === 独角兽相关融合 ===
	// 兔子 + 凤凰 = 独角兽
	registry.Register(pet.SpeciesRabbit, pet.SpeciesPhoenix, pet.HiddenSpeciesConfig{
		ResultSpecies:    pet.SpeciesUnicorn,
		TriggerThreshold: 200,
		Rarity:           5,
	})
	// 狗 + 凤凰 = 独角兽
	registry.Register(pet.SpeciesDog, pet.SpeciesPhoenix, pet.HiddenSpeciesConfig{
		ResultSpecies:    pet.SpeciesUnicorn,
		TriggerThreshold: 220,
		Rarity:           5,
	})

	// === 凤凰变种融合 ===
	// 猫头鹰 + 凤凰 = 更稀有的凤凰
	registry.Register(pet.SpeciesOwl, pet.SpeciesPhoenix, pet.HiddenSpeciesConfig{
		ResultSpecies:    pet.SpeciesPhoenix,
		TriggerThreshold: 180,
		Rarity:           5,
	})
	// 鹦鹉 + 火元素 = 凤凰
	registry.Register(pet.SpeciesParrot, pet.SpeciesFireSpirit, pet.HiddenSpeciesConfig{
		ResultSpecies:    pet.SpeciesPhoenix,
		TriggerThreshold: 190,
		Rarity:           5,
	})

	// === 史莱姆融合 ===
	// 史莱姆 + 水元素 = 水史莱姆变种
	registry.Register(pet.SpeciesSlime, pet.SpeciesWaterSpirit, pet.HiddenSpeciesConfig{
		ResultSpecies:    pet.SpeciesSlime, // 产生特殊变种
		TriggerThreshold: 150,
		Rarity:           3,
	})

	return registry
}

// registerMammals 注册哺乳类物种
func registerMammals(registry *pet.SpeciesRegistry) {
	// 猫
	registry.Register(&pet.Species{
		ID:       pet.SpeciesCat,
		Name:     "猫",
		Category: pet.CategoryMammal,
		BaseParts: []pet.PartType{
			pet.PartTypeNone, // 体色、体型等通用属性
		},
		SpecialParts: []pet.PartType{
			pet.PartTypeEar,
			pet.PartTypeTail,
			pet.PartTypeFur,
			pet.PartTypeWhisker,
		},
		Rarity:      1,
		IsHidden:    false,
		GenderRule:  pet.DefaultGenderRule(),
		BreedRules:  pet.DefaultBreedingRules(),
		Interpreter: NewFelineInterpreter(),
	})

	// 狗
	registry.Register(&pet.Species{
		ID:       pet.SpeciesDog,
		Name:     "狗",
		Category: pet.CategoryMammal,
		BaseParts: []pet.PartType{
			pet.PartTypeNone,
		},
		SpecialParts: []pet.PartType{
			pet.PartTypeEar,
			pet.PartTypeTail,
			pet.PartTypeFur,
			pet.PartTypeWhisker,
		},
		Rarity:      1,
		IsHidden:    false,
		GenderRule:  pet.DefaultGenderRule(),
		BreedRules:  pet.DefaultBreedingRules(),
		Interpreter: NewCanineInterpreter(),
	})

	// 兔
	registry.Register(&pet.Species{
		ID:       pet.SpeciesRabbit,
		Name:     "兔子",
		Category: pet.CategoryMammal,
		BaseParts: []pet.PartType{
			pet.PartTypeNone,
		},
		SpecialParts: []pet.PartType{
			pet.PartTypeEar,
			pet.PartTypeTail,
			pet.PartTypeFur,
			pet.PartTypeWhisker,
		},
		Rarity:      1,
		IsHidden:    false,
		GenderRule:  pet.DefaultGenderRule(),
		BreedRules:  pet.DefaultBreedingRules(),
		Interpreter: NewFelineInterpreter(), // 暂用猫科解释器
	})
}

// registerAvians 注册鸟类物种
func registerAvians(registry *pet.SpeciesRegistry) {
	// 鹦鹉
	registry.Register(&pet.Species{
		ID:       pet.SpeciesParrot,
		Name:     "鹦鹉",
		Category: pet.CategoryAvian,
		BaseParts: []pet.PartType{
			pet.PartTypeNone,
		},
		SpecialParts: []pet.PartType{
			pet.PartTypeWing,
			pet.PartTypeBeak,
			pet.PartTypeCrest,
			pet.PartTypeTail,
		},
		Rarity:      2,
		IsHidden:    false,
		GenderRule:  pet.DefaultGenderRule(),
		BreedRules:  pet.DefaultBreedingRules(),
		Interpreter: NewParrotInterpreter(),
	})

	// 猫头鹰
	registry.Register(&pet.Species{
		ID:       pet.SpeciesOwl,
		Name:     "猫头鹰",
		Category: pet.CategoryAvian,
		BaseParts: []pet.PartType{
			pet.PartTypeNone,
		},
		SpecialParts: []pet.PartType{
			pet.PartTypeWing,
			pet.PartTypeBeak,
			pet.PartTypeCrest,
			pet.PartTypeTail,
		},
		Rarity:      3,
		IsHidden:    false,
		GenderRule:  pet.DefaultGenderRule(),
		BreedRules:  pet.DefaultBreedingRules(),
		Interpreter: NewOwlInterpreter(),
	})
}

// registerAquatics 注册水生类物种
func registerAquatics(registry *pet.SpeciesRegistry) {
	// 金鱼
	registry.Register(&pet.Species{
		ID:       pet.SpeciesGoldfish,
		Name:     "金鱼",
		Category: pet.CategoryFish,
		BaseParts: []pet.PartType{
			pet.PartTypeNone,
		},
		SpecialParts: []pet.PartType{
			pet.PartTypeFin,
			pet.PartTypeScale,
			pet.PartTypeTailFin,
			pet.PartTypeWhisker,
		},
		Rarity:      1,
		IsHidden:    false,
		GenderRule:  pet.DefaultGenderRule(),
		BreedRules:  pet.DefaultBreedingRules(),
		Interpreter: NewGoldfishInterpreter(),
	})

	// 热带鱼
	registry.Register(&pet.Species{
		ID:       pet.SpeciesTropicalFish,
		Name:     "热带鱼",
		Category: pet.CategoryFish,
		BaseParts: []pet.PartType{
			pet.PartTypeNone,
		},
		SpecialParts: []pet.PartType{
			pet.PartTypeFin,
			pet.PartTypeScale,
			pet.PartTypeTailFin,
			pet.PartTypeWhisker,
		},
		Rarity:      2,
		IsHidden:    false,
		GenderRule:  pet.DefaultGenderRule(),
		BreedRules:  pet.DefaultBreedingRules(),
		Interpreter: NewTropicalFishInterpreter(),
	})
}

// registerFantasy 注册幻想类物种
func registerFantasy(registry *pet.SpeciesRegistry) {
	// 史莱姆（无性别，可分裂繁殖）
	registry.Register(&pet.Species{
		ID:       pet.SpeciesSlime,
		Name:     "史莱姆",
		Category: pet.CategoryFantasy,
		BaseParts: []pet.PartType{
			pet.PartTypeNone,
		},
		SpecialParts: []pet.PartType{
			pet.PartTypeTail,  // 体型形状
			pet.PartTypeScale, // 质感
			pet.PartTypeHorn,  // 核心
			pet.PartTypeAura,  // 光环
		},
		Rarity:      2,
		IsHidden:    false,
		GenderRule:  pet.AsexualGenderRule(), // 无性别
		BreedRules:  pet.DefaultBreedingRules(),
		Interpreter: NewSlimeInterpreter(),
	})

	// 凤凰（支持雌雄同体）
	registry.Register(&pet.Species{
		ID:       pet.SpeciesPhoenix,
		Name:     "凤凰",
		Category: pet.CategoryFantasy,
		BaseParts: []pet.PartType{
			pet.PartTypeNone,
		},
		SpecialParts: []pet.PartType{
			pet.PartTypeWing,
			pet.PartTypeCrest,
			pet.PartTypeTail,
			pet.PartTypeAura,
		},
		Rarity:      4,
		IsHidden:    false,
		GenderRule:  pet.MixedGenderRule(40, 40, 20), // 40%雄 40%雌 20%雌雄同体
		BreedRules:  pet.DefaultBreedingRules(),
		Interpreter: NewPhoenixInterpreter(),
	})

	// 龙（隐藏物种）
	registry.Register(&pet.Species{
		ID:       pet.SpeciesDragon,
		Name:     "龙",
		Category: pet.CategoryFantasy,
		BaseParts: []pet.PartType{
			pet.PartTypeNone,
		},
		SpecialParts: []pet.PartType{
			pet.PartTypeWing,
			pet.PartTypeHorn,
			pet.PartTypeArmor,
			pet.PartTypeTail,
		},
		Rarity:      5,
		IsHidden:    true, // 隐藏物种
		GenderRule:  pet.DefaultGenderRule(),
		BreedRules:  pet.DefaultBreedingRules(),
		Interpreter: NewDragonInterpreter(),
	})

	// 格里芬（隐藏物种）
	registry.Register(&pet.Species{
		ID:       pet.SpeciesGriffin,
		Name:     "格里芬",
		Category: pet.CategoryFantasy,
		BaseParts: []pet.PartType{
			pet.PartTypeNone,
		},
		SpecialParts: []pet.PartType{
			pet.PartTypeWing,
			pet.PartTypeEar,
			pet.PartTypeTail,
			pet.PartTypeClaw,
		},
		Rarity:      5,
		IsHidden:    true, // 隐藏物种
		GenderRule:  pet.DefaultGenderRule(),
		BreedRules:  pet.DefaultBreedingRules(),
		Interpreter: NewGriffinInterpreter(),
	})

	// 独角兽（隐藏物种）
	registry.Register(&pet.Species{
		ID:       pet.SpeciesUnicorn,
		Name:     "独角兽",
		Category: pet.CategoryFantasy,
		BaseParts: []pet.PartType{
			pet.PartTypeNone,
		},
		SpecialParts: []pet.PartType{
			pet.PartTypeHorn,
			pet.PartTypeFur,
			pet.PartTypeTail,
			pet.PartTypeAura,
		},
		Rarity:      5,
		IsHidden:    true, // 隐藏物种
		GenderRule:  pet.DefaultGenderRule(),
		BreedRules:  pet.DefaultBreedingRules(),
		Interpreter: NewUnicornInterpreter(),
	})

	// 火元素（无性别）
	registry.Register(&pet.Species{
		ID:       pet.SpeciesFireSpirit,
		Name:     "火元素",
		Category: pet.CategoryElemental,
		BaseParts: []pet.PartType{
			pet.PartTypeNone,
		},
		SpecialParts: []pet.PartType{
			pet.PartTypeAura,
		},
		Rarity:      3,
		IsHidden:    false,
		GenderRule:  pet.AsexualGenderRule(),
		BreedRules:  pet.DefaultBreedingRules(),
		Interpreter: NewSlimeInterpreter(), // 复用史莱姆解释器
	})

	// 水元素（无性别）
	registry.Register(&pet.Species{
		ID:       pet.SpeciesWaterSpirit,
		Name:     "水元素",
		Category: pet.CategoryElemental,
		BaseParts: []pet.PartType{
			pet.PartTypeNone,
		},
		SpecialParts: []pet.PartType{
			pet.PartTypeAura,
		},
		Rarity:      3,
		IsHidden:    false,
		GenderRule:  pet.AsexualGenderRule(),
		BreedRules:  pet.DefaultBreedingRules(),
		Interpreter: NewSlimeInterpreter(), // 复用史莱姆解释器
	})
}
