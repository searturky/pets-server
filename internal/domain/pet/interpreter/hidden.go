// Package interpreter 物种基因解释器
// 隐藏物种解释器 - 解释隐藏物种的特征基因
package interpreter

import "pets-server/internal/domain/pet"

// GriffinInterpreter 格里芬解释器（猫+鸟的融合）
type GriffinInterpreter struct{}

// NewGriffinInterpreter 创建格里芬解释器
func NewGriffinInterpreter() *GriffinInterpreter {
	return &GriffinInterpreter{}
}

// GetSpeciesID 返回此解释器对应的物种ID
func (g *GriffinInterpreter) GetSpeciesID() pet.SpeciesID {
	return pet.SpeciesGriffin
}

// 格里芬特征样式定义
var (
	griffinWingTypes = []string{
		"鹰翼", "猫头鹰翼", "天使翼", "金翼",
		"羽毛翼", "光翼", "暗翼", "彩虹翼",
		"风暴翼", "雷电翼", "火焰翼", "冰霜翼",
		"巨翼", "小翼", "隐形翼", "神圣翼",
	}

	griffinEarTypes = []string{
		"猫耳", "鹰羽耳", "尖耳", "圆耳",
		"精灵耳", "狮耳", "虎耳", "豹耳",
		"毛绒耳", "羽毛耳", "金属耳", "水晶耳",
		"神话耳", "古老耳", "皇家耳", "野性耳",
	}

	griffinTailTypes = []string{
		"狮尾", "猫尾", "鸟尾", "羽尾",
		"鹰尾", "凤凰尾", "蓬松尾", "细长尾",
		"双尾", "三尾", "环尾", "剑尾",
		"闪电尾", "火焰尾", "风暴尾", "皇家尾",
	}

	griffinClawTypes = []string{
		"鹰爪", "狮爪", "虎爪", "豹爪",
		"金爪", "银爪", "铜爪", "铁爪",
		"利爪", "巨爪", "隐形爪", "闪电爪",
		"火焰爪", "冰霜爪", "神圣爪", "暗影爪",
	}
)

// InterpretSpecialFeatures 解释格里芬特有特征
func (g *GriffinInterpreter) InterpretSpecialFeatures(gene pet.Gene) pet.SpecialAppearance {
	special := pet.NewSpecialAppearance()

	// 特征A: 翅膀
	wingValue := gene.SpecialA()
	wingMod := gene.SpecialModA()
	special.AddPart(pet.PartAppearance{
		PartType: pet.PartTypeWing,
		Style:    griffinWingTypes[wingValue%len(griffinWingTypes)],
		Value:    wingValue,
		Modifier: wingMod,
	})

	// 特征B: 耳朵
	earValue := gene.SpecialB()
	earMod := gene.SpecialModB()
	special.AddPart(pet.PartAppearance{
		PartType: pet.PartTypeEar,
		Style:    griffinEarTypes[earValue%len(griffinEarTypes)],
		Value:    earValue,
		Modifier: earMod,
	})

	// 特征C: 尾巴
	tailValue := gene.SpecialC()
	tailMod := gene.SpecialModC()
	special.AddPart(pet.PartAppearance{
		PartType: pet.PartTypeTail,
		Style:    griffinTailTypes[tailValue%len(griffinTailTypes)],
		Value:    tailValue,
		Modifier: tailMod,
	})

	// 特征D: 爪子
	clawValue := gene.SpecialD()
	clawMod := gene.SpecialModD()
	special.AddPart(pet.PartAppearance{
		PartType: pet.PartTypeClaw,
		Style:    griffinClawTypes[clawValue%len(griffinClawTypes)],
		Value:    clawValue,
		Modifier: clawMod,
	})

	return special
}

// GetFeatureNames 获取特征名称映射
func (g *GriffinInterpreter) GetFeatureNames() map[pet.PartType][]string {
	return map[pet.PartType][]string{
		pet.PartTypeWing: griffinWingTypes,
		pet.PartTypeEar:  griffinEarTypes,
		pet.PartTypeTail: griffinTailTypes,
		pet.PartTypeClaw: griffinClawTypes,
	}
}

// UnicornInterpreter 独角兽解释器
type UnicornInterpreter struct{}

// NewUnicornInterpreter 创建独角兽解释器
func NewUnicornInterpreter() *UnicornInterpreter {
	return &UnicornInterpreter{}
}

// GetSpeciesID 返回此解释器对应的物种ID
func (u *UnicornInterpreter) GetSpeciesID() pet.SpeciesID {
	return pet.SpeciesUnicorn
}

// 独角兽特征样式定义
var (
	unicornHornTypes = []string{
		"水晶角", "光明角", "彩虹角", "金角",
		"银角", "珍珠角", "钻石角", "月光角",
		"星辰角", "神圣角", "自然角", "精灵角",
		"螺旋角", "直角", "弯角", "分叉角",
	}

	unicornManeTypes = []string{
		"流光鬃", "彩虹鬃", "金色鬃", "银色鬃",
		"白色鬃", "星空鬃", "月光鬃", "花瓣鬃",
		"丝绸鬃", "波浪鬃", "飘逸鬃", "蓬松鬃",
		"长鬃", "短鬃", "编织鬃", "自然鬃",
	}

	unicornTailTypes = []string{
		"流星尾", "彩虹尾", "金色尾", "银色尾",
		"白色尾", "星空尾", "月光尾", "花瓣尾",
		"丝绸尾", "波浪尾", "飘逸尾", "蓬松尾",
		"长尾", "短尾", "编织尾", "自然尾",
	}

	unicornAuraTypes = []string{
		"月光光环", "星辰光环", "彩虹光环", "神圣光环",
		"自然光环", "精灵光环", "梦幻光环", "纯净光环",
		"希望光环", "爱情光环", "治愈光环", "祝福光环",
		"保护光环", "智慧光环", "勇气光环", "和平光环",
	}
)

// InterpretSpecialFeatures 解释独角兽特有特征
func (u *UnicornInterpreter) InterpretSpecialFeatures(gene pet.Gene) pet.SpecialAppearance {
	special := pet.NewSpecialAppearance()

	// 特征A: 角
	hornValue := gene.SpecialA()
	hornMod := gene.SpecialModA()
	special.AddPart(pet.PartAppearance{
		PartType: pet.PartTypeHorn,
		Style:    unicornHornTypes[hornValue%len(unicornHornTypes)],
		Value:    hornValue,
		Modifier: hornMod,
	})

	// 特征B: 鬃毛（复用Fur类型）
	maneValue := gene.SpecialB()
	maneMod := gene.SpecialModB()
	special.AddPart(pet.PartAppearance{
		PartType: pet.PartTypeFur,
		Style:    unicornManeTypes[maneValue%len(unicornManeTypes)],
		Value:    maneValue,
		Modifier: maneMod,
	})

	// 特征C: 尾巴
	tailValue := gene.SpecialC()
	tailMod := gene.SpecialModC()
	special.AddPart(pet.PartAppearance{
		PartType: pet.PartTypeTail,
		Style:    unicornTailTypes[tailValue%len(unicornTailTypes)],
		Value:    tailValue,
		Modifier: tailMod,
	})

	// 特征D: 光环
	auraValue := gene.SpecialD()
	auraMod := gene.SpecialModD()
	special.AddPart(pet.PartAppearance{
		PartType: pet.PartTypeAura,
		Style:    unicornAuraTypes[auraValue%len(unicornAuraTypes)],
		Value:    auraValue,
		Modifier: auraMod,
	})

	return special
}

// GetFeatureNames 获取特征名称映射
func (u *UnicornInterpreter) GetFeatureNames() map[pet.PartType][]string {
	return map[pet.PartType][]string{
		pet.PartTypeHorn: unicornHornTypes,
		pet.PartTypeFur:  unicornManeTypes,
		pet.PartTypeTail: unicornTailTypes,
		pet.PartTypeAura: unicornAuraTypes,
	}
}
