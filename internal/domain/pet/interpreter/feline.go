// Package interpreter 物种基因解释器
// 猫科解释器 - 解释猫科动物的特征基因
package interpreter

import "pets-server/internal/domain/pet"

// FelineInterpreter 猫科解释器
type FelineInterpreter struct{}

// NewFelineInterpreter 创建猫科解释器
func NewFelineInterpreter() *FelineInterpreter {
	return &FelineInterpreter{}
}

// GetSpeciesID 返回此解释器对应的物种ID
func (f *FelineInterpreter) GetSpeciesID() pet.SpeciesID {
	return pet.SpeciesCat
}

// 猫科特征样式定义
var (
	felineEarTypes = []string{
		"立耳", "折耳", "卷耳", "圆耳",
		"尖耳", "垂耳", "小立耳", "大圆耳",
		"三角耳", "蝙蝠耳", "精灵耳", "毛耳",
		"短耳", "长耳", "宽耳", "窄耳",
	}

	felineTailTypes = []string{
		"长尾", "短尾", "蓬松尾", "细长尾",
		"卷尾", "直尾", "粗尾", "细尾",
		"狮尾", "松鼠尾", "毛球尾", "无尾",
		"弯尾", "环尾", "羽毛尾", "渐变尾",
	}

	felineFurTypes = []string{
		"短毛", "长毛", "卷毛", "波浪毛",
		"丝绒毛", "绒毛", "粗毛", "细毛",
		"双层毛", "单层毛", "蓬松毛", "贴身毛",
		"斑纹毛", "虎纹毛", "豹纹毛", "渐层毛",
	}

	felineWhiskerTypes = []string{
		"长胡须", "短胡须", "卷胡须", "直胡须",
		"粗胡须", "细胡须", "密胡须", "疏胡须",
		"白胡须", "黑胡须", "灰胡须", "彩虹胡须",
		"弯曲胡须", "蓬松胡须", "精致胡须", "威严胡须",
	}
)

// InterpretSpecialFeatures 解释猫科特有特征
func (f *FelineInterpreter) InterpretSpecialFeatures(gene pet.Gene) pet.SpecialAppearance {
	special := pet.NewSpecialAppearance()

	// 特征A: 耳朵
	earValue := gene.SpecialA()
	earMod := gene.SpecialModA()
	special.AddPart(pet.PartAppearance{
		PartType: pet.PartTypeEar,
		Style:    felineEarTypes[earValue%len(felineEarTypes)],
		Value:    earValue,
		Modifier: earMod,
	})

	// 特征B: 尾巴
	tailValue := gene.SpecialB()
	tailMod := gene.SpecialModB()
	special.AddPart(pet.PartAppearance{
		PartType: pet.PartTypeTail,
		Style:    felineTailTypes[tailValue%len(felineTailTypes)],
		Value:    tailValue,
		Modifier: tailMod,
	})

	// 特征C: 毛纹
	furValue := gene.SpecialC()
	furMod := gene.SpecialModC()
	special.AddPart(pet.PartAppearance{
		PartType: pet.PartTypeFur,
		Style:    felineFurTypes[furValue%len(felineFurTypes)],
		Value:    furValue,
		Modifier: furMod,
	})

	// 特征D: 胡须
	whiskerValue := gene.SpecialD()
	whiskerMod := gene.SpecialModD()
	special.AddPart(pet.PartAppearance{
		PartType: pet.PartTypeWhisker,
		Style:    felineWhiskerTypes[whiskerValue%len(felineWhiskerTypes)],
		Value:    whiskerValue,
		Modifier: whiskerMod,
	})

	return special
}

// GetFeatureNames 获取特征名称映射
func (f *FelineInterpreter) GetFeatureNames() map[pet.PartType][]string {
	return map[pet.PartType][]string{
		pet.PartTypeEar:     felineEarTypes,
		pet.PartTypeTail:    felineTailTypes,
		pet.PartTypeFur:     felineFurTypes,
		pet.PartTypeWhisker: felineWhiskerTypes,
	}
}

// CanineInterpreter 犬科解释器
type CanineInterpreter struct{}

// NewCanineInterpreter 创建犬科解释器
func NewCanineInterpreter() *CanineInterpreter {
	return &CanineInterpreter{}
}

// GetSpeciesID 返回此解释器对应的物种ID
func (c *CanineInterpreter) GetSpeciesID() pet.SpeciesID {
	return pet.SpeciesDog
}

// 犬科特征样式定义
var (
	canineEarTypes = []string{
		"立耳", "垂耳", "半立耳", "玫瑰耳",
		"蝙蝠耳", "纽扣耳", "折耳", "飞耳",
		"三角耳", "圆耳", "尖耳", "大耳",
		"小耳", "毛耳", "薄耳", "厚耳",
	}

	canineTailTypes = []string{
		"卷尾", "直尾", "弯尾", "镰刀尾",
		"螺旋尾", "断尾", "扫帚尾", "剑尾",
		"羽毛尾", "松鼠尾", "鞭尾", "环尾",
		"低垂尾", "高举尾", "蓬松尾", "光滑尾",
	}

	canineFurTypes = []string{
		"短毛", "长毛", "卷毛", "钢丝毛",
		"丝毛", "绒毛", "双层毛", "单层毛",
		"蓬松毛", "光滑毛", "粗毛", "细毛",
		"波浪毛", "直毛", "杂毛", "缎子毛",
	}

	canineMuzzleTypes = []string{
		"长吻", "短吻", "方吻", "尖吻",
		"宽吻", "窄吻", "扁吻", "翘吻",
		"黑鼻", "粉鼻", "斑点鼻", "蝴蝶鼻",
		"大鼻", "小鼻", "湿鼻", "干鼻",
	}
)

// InterpretSpecialFeatures 解释犬科特有特征
func (c *CanineInterpreter) InterpretSpecialFeatures(gene pet.Gene) pet.SpecialAppearance {
	special := pet.NewSpecialAppearance()

	// 特征A: 耳朵
	earValue := gene.SpecialA()
	earMod := gene.SpecialModA()
	special.AddPart(pet.PartAppearance{
		PartType: pet.PartTypeEar,
		Style:    canineEarTypes[earValue%len(canineEarTypes)],
		Value:    earValue,
		Modifier: earMod,
	})

	// 特征B: 尾巴
	tailValue := gene.SpecialB()
	tailMod := gene.SpecialModB()
	special.AddPart(pet.PartAppearance{
		PartType: pet.PartTypeTail,
		Style:    canineTailTypes[tailValue%len(canineTailTypes)],
		Value:    tailValue,
		Modifier: tailMod,
	})

	// 特征C: 毛纹
	furValue := gene.SpecialC()
	furMod := gene.SpecialModC()
	special.AddPart(pet.PartAppearance{
		PartType: pet.PartTypeFur,
		Style:    canineFurTypes[furValue%len(canineFurTypes)],
		Value:    furValue,
		Modifier: furMod,
	})

	// 特征D: 吻部（复用Whisker类型）
	muzzleValue := gene.SpecialD()
	muzzleMod := gene.SpecialModD()
	special.AddPart(pet.PartAppearance{
		PartType: pet.PartTypeWhisker, // 复用胡须类型表示吻部
		Style:    canineMuzzleTypes[muzzleValue%len(canineMuzzleTypes)],
		Value:    muzzleValue,
		Modifier: muzzleMod,
	})

	return special
}

// GetFeatureNames 获取特征名称映射
func (c *CanineInterpreter) GetFeatureNames() map[pet.PartType][]string {
	return map[pet.PartType][]string{
		pet.PartTypeEar:     canineEarTypes,
		pet.PartTypeTail:    canineTailTypes,
		pet.PartTypeFur:     canineFurTypes,
		pet.PartTypeWhisker: canineMuzzleTypes,
	}
}
