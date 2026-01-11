// Package interpreter 物种基因解释器
// 鸟类解释器 - 解释鸟类动物的特征基因
package interpreter

import "pets-server/internal/domain/pet"

// AvianInterpreter 鸟类解释器
type AvianInterpreter struct {
	speciesID pet.SpeciesID
}

// NewAvianInterpreter 创建鸟类解释器
func NewAvianInterpreter(speciesID pet.SpeciesID) *AvianInterpreter {
	return &AvianInterpreter{speciesID: speciesID}
}

// NewParrotInterpreter 创建鹦鹉解释器
func NewParrotInterpreter() *AvianInterpreter {
	return NewAvianInterpreter(pet.SpeciesParrot)
}

// NewOwlInterpreter 创建猫头鹰解释器
func NewOwlInterpreter() *AvianInterpreter {
	return NewAvianInterpreter(pet.SpeciesOwl)
}

// GetSpeciesID 返回此解释器对应的物种ID
func (a *AvianInterpreter) GetSpeciesID() pet.SpeciesID {
	return a.speciesID
}

// 鸟类特征样式定义
var (
	avianWingTypes = []string{
		"小翅膀", "圆翅膀", "尖翅膀", "大翅膀",
		"羽翼", "天使翼", "蝙蝠翼", "蜻蜓翼",
		"宽翅膀", "窄翅膀", "彩虹翅膀", "透明翅膀",
		"羽毛翅膀", "绒毛翅膀", "光翅膀", "暗翅膀",
	}

	avianBeakTypes = []string{
		"短喙", "长喙", "弯喙", "直喙",
		"尖喙", "钝喙", "宽喙", "窄喙",
		"勾喙", "扁喙", "锥形喙", "镊子喙",
		"彩色喙", "黑喙", "橙喙", "红喙",
	}

	avianCrestTypes = []string{
		"无羽冠", "小羽冠", "大羽冠", "竖羽冠",
		"扇形羽冠", "刺羽冠", "冠羽", "凤冠",
		"莫西干羽冠", "卷羽冠", "蓬松羽冠", "流线羽冠",
		"彩虹羽冠", "渐变羽冠", "发光羽冠", "皇冠羽冠",
	}

	avianTailTypes = []string{
		"短尾羽", "长尾羽", "扇形尾", "剪尾",
		"燕尾", "圆尾", "尖尾", "方尾",
		"孔雀尾", "凤凰尾", "彩带尾", "渐变尾",
		"蓬松尾", "细长尾", "卷尾", "分叉尾",
	}
)

// InterpretSpecialFeatures 解释鸟类特有特征
func (a *AvianInterpreter) InterpretSpecialFeatures(gene pet.Gene) pet.SpecialAppearance {
	special := pet.NewSpecialAppearance()

	// 特征A: 翅膀
	wingValue := gene.SpecialA()
	wingMod := gene.SpecialModA()
	special.AddPart(pet.PartAppearance{
		PartType: pet.PartTypeWing,
		Style:    avianWingTypes[wingValue%len(avianWingTypes)],
		Value:    wingValue,
		Modifier: wingMod,
	})

	// 特征B: 喙
	beakValue := gene.SpecialB()
	beakMod := gene.SpecialModB()
	special.AddPart(pet.PartAppearance{
		PartType: pet.PartTypeBeak,
		Style:    avianBeakTypes[beakValue%len(avianBeakTypes)],
		Value:    beakValue,
		Modifier: beakMod,
	})

	// 特征C: 羽冠
	crestValue := gene.SpecialC()
	crestMod := gene.SpecialModC()
	special.AddPart(pet.PartAppearance{
		PartType: pet.PartTypeCrest,
		Style:    avianCrestTypes[crestValue%len(avianCrestTypes)],
		Value:    crestValue,
		Modifier: crestMod,
	})

	// 特征D: 尾羽
	tailValue := gene.SpecialD()
	tailMod := gene.SpecialModD()
	special.AddPart(pet.PartAppearance{
		PartType: pet.PartTypeTail,
		Style:    avianTailTypes[tailValue%len(avianTailTypes)],
		Value:    tailValue,
		Modifier: tailMod,
	})

	return special
}

// GetFeatureNames 获取特征名称映射
func (a *AvianInterpreter) GetFeatureNames() map[pet.PartType][]string {
	return map[pet.PartType][]string{
		pet.PartTypeWing:  avianWingTypes,
		pet.PartTypeBeak:  avianBeakTypes,
		pet.PartTypeCrest: avianCrestTypes,
		pet.PartTypeTail:  avianTailTypes,
	}
}
