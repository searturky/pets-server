// Package interpreter 物种基因解释器
// 幻想类解释器 - 解释幻想生物的特征基因
package interpreter

import "pets-server/internal/domain/pet"

// DragonInterpreter 龙类解释器
type DragonInterpreter struct{}

// NewDragonInterpreter 创建龙类解释器
func NewDragonInterpreter() *DragonInterpreter {
	return &DragonInterpreter{}
}

// GetSpeciesID 返回此解释器对应的物种ID
func (d *DragonInterpreter) GetSpeciesID() pet.SpeciesID {
	return pet.SpeciesDragon
}

// 龙类特征样式定义
var (
	dragonWingTypes = []string{
		"蝙蝠翼", "羽翼", "骨翼", "膜翼",
		"水晶翼", "火焰翼", "冰霜翼", "雷电翼",
		"暗影翼", "光明翼", "双翼", "四翼",
		"小翼", "巨翼", "破损翼", "完美翼",
	}

	dragonHornTypes = []string{
		"单角", "双角", "弯角", "直角",
		"螺旋角", "分叉角", "皇冠角", "鹿角",
		"水晶角", "火焰角", "冰霜角", "雷电角",
		"暗影角", "光明角", "古老角", "神圣角",
	}

	dragonArmorTypes = []string{
		"细鳞甲", "厚鳞甲", "板甲", "骨甲",
		"水晶甲", "火岩甲", "冰晶甲", "雷纹甲",
		"暗影甲", "光明甲", "古老甲", "神圣甲",
		"熔岩甲", "海洋甲", "森林甲", "天空甲",
	}

	dragonTailTypes = []string{
		"尖刺尾", "锤尾", "剑尾", "羽尾",
		"蛇尾", "鞭尾", "火焰尾", "冰霜尾",
		"雷电尾", "毒尾", "光尾", "暗尾",
		"分叉尾", "环尾", "锯齿尾", "柔软尾",
	}
)

// InterpretSpecialFeatures 解释龙类特有特征
func (d *DragonInterpreter) InterpretSpecialFeatures(gene pet.Gene) pet.SpecialAppearance {
	special := pet.NewSpecialAppearance()

	// 特征A: 翅膀
	wingValue := gene.SpecialA()
	wingMod := gene.SpecialModA()
	special.AddPart(pet.PartAppearance{
		PartType: pet.PartTypeWing,
		Style:    dragonWingTypes[wingValue%len(dragonWingTypes)],
		Value:    wingValue,
		Modifier: wingMod,
	})

	// 特征B: 角
	hornValue := gene.SpecialB()
	hornMod := gene.SpecialModB()
	special.AddPart(pet.PartAppearance{
		PartType: pet.PartTypeHorn,
		Style:    dragonHornTypes[hornValue%len(dragonHornTypes)],
		Value:    hornValue,
		Modifier: hornMod,
	})

	// 特征C: 鳞甲
	armorValue := gene.SpecialC()
	armorMod := gene.SpecialModC()
	special.AddPart(pet.PartAppearance{
		PartType: pet.PartTypeArmor,
		Style:    dragonArmorTypes[armorValue%len(dragonArmorTypes)],
		Value:    armorValue,
		Modifier: armorMod,
	})

	// 特征D: 尾巴
	tailValue := gene.SpecialD()
	tailMod := gene.SpecialModD()
	special.AddPart(pet.PartAppearance{
		PartType: pet.PartTypeTail,
		Style:    dragonTailTypes[tailValue%len(dragonTailTypes)],
		Value:    tailValue,
		Modifier: tailMod,
	})

	return special
}

// GetFeatureNames 获取特征名称映射
func (d *DragonInterpreter) GetFeatureNames() map[pet.PartType][]string {
	return map[pet.PartType][]string{
		pet.PartTypeWing:  dragonWingTypes,
		pet.PartTypeHorn:  dragonHornTypes,
		pet.PartTypeArmor: dragonArmorTypes,
		pet.PartTypeTail:  dragonTailTypes,
	}
}

// SlimeInterpreter 史莱姆解释器
type SlimeInterpreter struct{}

// NewSlimeInterpreter 创建史莱姆解释器
func NewSlimeInterpreter() *SlimeInterpreter {
	return &SlimeInterpreter{}
}

// GetSpeciesID 返回此解释器对应的物种ID
func (s *SlimeInterpreter) GetSpeciesID() pet.SpeciesID {
	return pet.SpeciesSlime
}

// 史莱姆特征样式定义
var (
	slimeBodyTypes = []string{
		"圆形", "椭圆形", "水滴形", "不规则形",
		"星形", "心形", "方形", "三角形",
		"云朵形", "气泡形", "果冻形", "布丁形",
		"融化形", "弹跳形", "分裂形", "合体形",
	}

	slimeTextureTypes = []string{
		"透明", "半透明", "不透明", "发光",
		"闪烁", "渐变", "彩虹", "金属",
		"果冻质感", "水晶质感", "凝胶质感", "液态质感",
		"气泡质感", "星空质感", "云朵质感", "熔岩质感",
	}

	slimeCoreTypes = []string{
		"无核", "单核", "多核", "星核",
		"心核", "水晶核", "火焰核", "冰霜核",
		"雷电核", "暗核", "光核", "彩虹核",
		"漂浮核", "旋转核", "脉动核", "神秘核",
	}

	slimeAuraTypes = []string{
		"无光环", "淡光环", "强光环", "脉动光环",
		"彩虹光环", "火焰光环", "冰霜光环", "雷电光环",
		"暗影光环", "神圣光环", "自然光环", "星空光环",
		"水波光环", "气泡光环", "花瓣光环", "雪花光环",
	}
)

// InterpretSpecialFeatures 解释史莱姆特有特征
func (s *SlimeInterpreter) InterpretSpecialFeatures(gene pet.Gene) pet.SpecialAppearance {
	special := pet.NewSpecialAppearance()

	// 特征A: 体型形状（复用Tail类型）
	bodyValue := gene.SpecialA()
	bodyMod := gene.SpecialModA()
	special.AddPart(pet.PartAppearance{
		PartType: pet.PartTypeTail, // 复用表示体型形状
		Style:    slimeBodyTypes[bodyValue%len(slimeBodyTypes)],
		Value:    bodyValue,
		Modifier: bodyMod,
	})

	// 特征B: 质感（复用Scale类型）
	textureValue := gene.SpecialB()
	textureMod := gene.SpecialModB()
	special.AddPart(pet.PartAppearance{
		PartType: pet.PartTypeScale, // 复用表示质感
		Style:    slimeTextureTypes[textureValue%len(slimeTextureTypes)],
		Value:    textureValue,
		Modifier: textureMod,
	})

	// 特征C: 核心（复用Horn类型）
	coreValue := gene.SpecialC()
	coreMod := gene.SpecialModC()
	special.AddPart(pet.PartAppearance{
		PartType: pet.PartTypeHorn, // 复用表示核心
		Style:    slimeCoreTypes[coreValue%len(slimeCoreTypes)],
		Value:    coreValue,
		Modifier: coreMod,
	})

	// 特征D: 光环
	auraValue := gene.SpecialD()
	auraMod := gene.SpecialModD()
	special.AddPart(pet.PartAppearance{
		PartType: pet.PartTypeAura,
		Style:    slimeAuraTypes[auraValue%len(slimeAuraTypes)],
		Value:    auraValue,
		Modifier: auraMod,
	})

	return special
}

// GetFeatureNames 获取特征名称映射
func (s *SlimeInterpreter) GetFeatureNames() map[pet.PartType][]string {
	return map[pet.PartType][]string{
		pet.PartTypeTail:  slimeBodyTypes,
		pet.PartTypeScale: slimeTextureTypes,
		pet.PartTypeHorn:  slimeCoreTypes,
		pet.PartTypeAura:  slimeAuraTypes,
	}
}

// PhoenixInterpreter 凤凰解释器
type PhoenixInterpreter struct{}

// NewPhoenixInterpreter 创建凤凰解释器
func NewPhoenixInterpreter() *PhoenixInterpreter {
	return &PhoenixInterpreter{}
}

// GetSpeciesID 返回此解释器对应的物种ID
func (p *PhoenixInterpreter) GetSpeciesID() pet.SpeciesID {
	return pet.SpeciesPhoenix
}

// 凤凰特征样式定义
var (
	phoenixWingTypes = []string{
		"火焰翼", "光明翼", "彩虹翼", "金翼",
		"红翼", "橙翼", "紫翼", "白翼",
		"凤凰翼", "朱雀翼", "涅槃翼", "永恒翼",
		"烈焰翼", "神圣翼", "太阳翼", "星辰翼",
	}

	phoenixCrestTypes = []string{
		"火焰冠", "光明冠", "彩虹冠", "金冠",
		"羽冠", "凤冠", "皇冠", "神冠",
		"太阳冠", "星辰冠", "流光冠", "永恒冠",
		"朱雀冠", "涅槃冠", "重生冠", "神圣冠",
	}

	phoenixTailTypes = []string{
		"火焰尾", "光明尾", "彩虹尾", "金尾",
		"孔雀尾", "凤凰尾", "流光尾", "星辰尾",
		"九尾", "三尾", "七尾", "单尾",
		"朱雀尾", "涅槃尾", "重生尾", "永恒尾",
	}

	phoenixAuraTypes = []string{
		"火焰光环", "光明光环", "彩虹光环", "金色光环",
		"太阳光环", "星辰光环", "凤凰光环", "涅槃光环",
		"重生光环", "永恒光环", "神圣光环", "灼热光环",
		"温暖光环", "希望光环", "生命光环", "奇迹光环",
	}
)

// InterpretSpecialFeatures 解释凤凰特有特征
func (p *PhoenixInterpreter) InterpretSpecialFeatures(gene pet.Gene) pet.SpecialAppearance {
	special := pet.NewSpecialAppearance()

	// 特征A: 翅膀
	wingValue := gene.SpecialA()
	wingMod := gene.SpecialModA()
	special.AddPart(pet.PartAppearance{
		PartType: pet.PartTypeWing,
		Style:    phoenixWingTypes[wingValue%len(phoenixWingTypes)],
		Value:    wingValue,
		Modifier: wingMod,
	})

	// 特征B: 羽冠
	crestValue := gene.SpecialB()
	crestMod := gene.SpecialModB()
	special.AddPart(pet.PartAppearance{
		PartType: pet.PartTypeCrest,
		Style:    phoenixCrestTypes[crestValue%len(phoenixCrestTypes)],
		Value:    crestValue,
		Modifier: crestMod,
	})

	// 特征C: 尾羽
	tailValue := gene.SpecialC()
	tailMod := gene.SpecialModC()
	special.AddPart(pet.PartAppearance{
		PartType: pet.PartTypeTail,
		Style:    phoenixTailTypes[tailValue%len(phoenixTailTypes)],
		Value:    tailValue,
		Modifier: tailMod,
	})

	// 特征D: 光环
	auraValue := gene.SpecialD()
	auraMod := gene.SpecialModD()
	special.AddPart(pet.PartAppearance{
		PartType: pet.PartTypeAura,
		Style:    phoenixAuraTypes[auraValue%len(phoenixAuraTypes)],
		Value:    auraValue,
		Modifier: auraMod,
	})

	return special
}

// GetFeatureNames 获取特征名称映射
func (p *PhoenixInterpreter) GetFeatureNames() map[pet.PartType][]string {
	return map[pet.PartType][]string{
		pet.PartTypeWing:  phoenixWingTypes,
		pet.PartTypeCrest: phoenixCrestTypes,
		pet.PartTypeTail:  phoenixTailTypes,
		pet.PartTypeAura:  phoenixAuraTypes,
	}
}
