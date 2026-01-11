// Package interpreter 物种基因解释器
// 水生类解释器 - 解释鱼类动物的特征基因
package interpreter

import "pets-server/internal/domain/pet"

// AquaticInterpreter 水生类解释器
type AquaticInterpreter struct {
	speciesID pet.SpeciesID
}

// NewAquaticInterpreter 创建水生类解释器
func NewAquaticInterpreter(speciesID pet.SpeciesID) *AquaticInterpreter {
	return &AquaticInterpreter{speciesID: speciesID}
}

// NewGoldfishInterpreter 创建金鱼解释器
func NewGoldfishInterpreter() *AquaticInterpreter {
	return NewAquaticInterpreter(pet.SpeciesGoldfish)
}

// NewTropicalFishInterpreter 创建热带鱼解释器
func NewTropicalFishInterpreter() *AquaticInterpreter {
	return NewAquaticInterpreter(pet.SpeciesTropicalFish)
}

// GetSpeciesID 返回此解释器对应的物种ID
func (a *AquaticInterpreter) GetSpeciesID() pet.SpeciesID {
	return a.speciesID
}

// 水生类特征样式定义
var (
	aquaticFinTypes = []string{
		"小鳍", "大鳍", "扇形鳍", "三角鳍",
		"流线鳍", "羽毛鳍", "透明鳍", "彩色鳍",
		"尖鳍", "圆鳍", "锯齿鳍", "飘逸鳍",
		"硬鳍", "软鳍", "分叉鳍", "连续鳍",
	}

	aquaticScaleTypes = []string{
		"细鳞", "大鳞", "无鳞", "菱形鳞",
		"圆鳞", "栉鳞", "盾鳞", "骨板",
		"闪光鳞", "哑光鳞", "彩虹鳞", "透明鳞",
		"金属鳞", "珍珠鳞", "渐变鳞", "斑点鳞",
	}

	aquaticTailFinTypes = []string{
		"扇尾", "剪刀尾", "燕尾", "圆尾",
		"尖尾", "琴尾", "蝶尾", "凤尾",
		"双尾", "三尾", "裙尾", "飘带尾",
		"狮子尾", "孔雀尾", "流星尾", "彗星尾",
	}

	aquaticWhiskerTypes = []string{
		"无须", "短须", "长须", "卷须",
		"双须", "四须", "六须", "八须",
		"粗须", "细须", "透明须", "发光须",
		"触手须", "羽毛须", "分叉须", "蓬松须",
	}
)

// InterpretSpecialFeatures 解释水生类特有特征
func (a *AquaticInterpreter) InterpretSpecialFeatures(gene pet.Gene) pet.SpecialAppearance {
	special := pet.NewSpecialAppearance()

	// 特征A: 背鳍
	finValue := gene.SpecialA()
	finMod := gene.SpecialModA()
	special.AddPart(pet.PartAppearance{
		PartType: pet.PartTypeFin,
		Style:    aquaticFinTypes[finValue%len(aquaticFinTypes)],
		Value:    finValue,
		Modifier: finMod,
	})

	// 特征B: 鳞片
	scaleValue := gene.SpecialB()
	scaleMod := gene.SpecialModB()
	special.AddPart(pet.PartAppearance{
		PartType: pet.PartTypeScale,
		Style:    aquaticScaleTypes[scaleValue%len(aquaticScaleTypes)],
		Value:    scaleValue,
		Modifier: scaleMod,
	})

	// 特征C: 尾鳍
	tailFinValue := gene.SpecialC()
	tailFinMod := gene.SpecialModC()
	special.AddPart(pet.PartAppearance{
		PartType: pet.PartTypeTailFin,
		Style:    aquaticTailFinTypes[tailFinValue%len(aquaticTailFinTypes)],
		Value:    tailFinValue,
		Modifier: tailFinMod,
	})

	// 特征D: 触须
	whiskerValue := gene.SpecialD()
	whiskerMod := gene.SpecialModD()
	special.AddPart(pet.PartAppearance{
		PartType: pet.PartTypeWhisker,
		Style:    aquaticWhiskerTypes[whiskerValue%len(aquaticWhiskerTypes)],
		Value:    whiskerValue,
		Modifier: whiskerMod,
	})

	return special
}

// GetFeatureNames 获取特征名称映射
func (a *AquaticInterpreter) GetFeatureNames() map[pet.PartType][]string {
	return map[pet.PartType][]string{
		pet.PartTypeFin:     aquaticFinTypes,
		pet.PartTypeScale:   aquaticScaleTypes,
		pet.PartTypeTailFin: aquaticTailFinTypes,
		pet.PartTypeWhisker: aquaticWhiskerTypes,
	}
}
