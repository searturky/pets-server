// Package pet 宠物领域
// Gender 性别系统 - 定义性别类型和繁衍配对规则
package pet

// Gender 性别类型
type Gender int

const (
	GenderNone          Gender = 0 // 无性别（元素生物、史莱姆等）
	GenderMale          Gender = 1 // 雄性
	GenderFemale        Gender = 2 // 雌性
	GenderHermaphrodite Gender = 3 // 雌雄同体（蜗牛、某些幻想生物）
)

// Name 性别名称
func (g Gender) Name() string {
	names := map[Gender]string{
		GenderNone:          "无",
		GenderMale:          "雄性",
		GenderFemale:        "雌性",
		GenderHermaphrodite: "雌雄同体",
	}
	if name, ok := names[g]; ok {
		return name
	}
	return "未知"
}

// Symbol 性别符号
func (g Gender) Symbol() string {
	symbols := map[Gender]string{
		GenderNone:          "○",
		GenderMale:          "♂",
		GenderFemale:        "♀",
		GenderHermaphrodite: "⚥",
	}
	if symbol, ok := symbols[g]; ok {
		return symbol
	}
	return "?"
}

// IsValid 是否为有效性别
func (g Gender) IsValid() bool {
	return g >= GenderNone && g <= GenderHermaphrodite
}

// CanBreed 是否可以参与繁殖
func (g Gender) CanBreed() bool {
	return g != GenderNone
}

// CanSelfBreed 是否可以自我繁殖
func (g Gender) CanSelfBreed() bool {
	return g == GenderNone
}

// CanBreedWith 判断两个性别是否可以配对繁衍
func CanBreedWith(g1, g2 Gender) bool {
	// 无性别不能与他人配对，只能自我繁殖
	if g1 == GenderNone || g2 == GenderNone {
		return false
	}

	compatibilityMatrix := map[Gender][]Gender{
		GenderMale:          {GenderFemale, GenderHermaphrodite},
		GenderFemale:        {GenderMale, GenderHermaphrodite},
		GenderHermaphrodite: {GenderMale, GenderFemale, GenderHermaphrodite},
	}

	compatibleGenders, ok := compatibilityMatrix[g1]
	if !ok {
		return false
	}

	for _, compatible := range compatibleGenders {
		if compatible == g2 {
			return true
		}
	}
	return false
}

// DetermineGender 根据基因和物种规则确定性别
func DetermineGender(gene Gene, rule GenderRule) Gender {
	// 如果物种只有一种性别
	if len(rule.AllowedGenders) == 1 {
		return rule.AllowedGenders[0]
	}

	// 使用基因位置0的值来决定性别
	genderValue := gene.HexAt(0)

	// 计算总比例
	total := 0
	for _, ratio := range rule.DefaultRatio {
		total += ratio
	}

	if total == 0 {
		return rule.AllowedGenders[0]
	}

	// 将基因值映射到比例范围
	threshold := genderValue * total / 16

	// 按比例分配
	accumulated := 0
	for _, gender := range rule.AllowedGenders {
		ratio, ok := rule.DefaultRatio[gender]
		if !ok {
			continue
		}
		accumulated += ratio
		if threshold < accumulated {
			return gender
		}
	}

	return rule.AllowedGenders[0]
}

// DetermineChildGender 确定子代性别
// 考虑父母性别对子代的影响
func DetermineChildGender(parent1, parent2 *Pet, childGene Gene, rule GenderRule) Gender {
	// 如果物种只有一种性别
	if len(rule.AllowedGenders) == 1 {
		return rule.AllowedGenders[0]
	}

	// 双亲都是雌雄同体时，子代有更高概率也是雌雄同体
	if parent1 != nil && parent2 != nil {
		if parent1.Gender == GenderHermaphrodite && parent2.Gender == GenderHermaphrodite {
			// 检查是否允许雌雄同体
			for _, g := range rule.AllowedGenders {
				if g == GenderHermaphrodite {
					// 60% 概率继承雌雄同体
					if randomInt(100) < 60 {
						return GenderHermaphrodite
					}
					break
				}
			}
		}
	}

	// 使用常规规则确定性别
	return DetermineGender(childGene, rule)
}

// GenderAppearanceModifier 性别外观修饰
// 某些物种的雄性和雌性在外观上有差异
type GenderAppearanceModifier struct {
	Gender         Gender
	ColorModifier  float64 // 颜色饱和度修饰 (1.0 = 不变, >1.0 = 更鲜艳)
	SizeModifier   float64 // 体型修饰 (1.0 = 不变)
	SpecialFeature string  // 特殊特征描述
}

// DefaultGenderModifiers 默认性别外观修饰
var DefaultGenderModifiers = map[Gender]GenderAppearanceModifier{
	GenderNone: {
		Gender:        GenderNone,
		ColorModifier: 1.0,
		SizeModifier:  1.0,
	},
	GenderMale: {
		Gender:        GenderMale,
		ColorModifier: 1.1, // 雄性通常颜色更鲜艳
		SizeModifier:  1.05,
	},
	GenderFemale: {
		Gender:        GenderFemale,
		ColorModifier: 1.0,
		SizeModifier:  0.95,
	},
	GenderHermaphrodite: {
		Gender:        GenderHermaphrodite,
		ColorModifier: 1.05,
		SizeModifier:  1.0,
	},
}

// GetGenderModifier 获取性别外观修饰
func GetGenderModifier(gender Gender) GenderAppearanceModifier {
	if modifier, ok := DefaultGenderModifiers[gender]; ok {
		return modifier
	}
	return GenderAppearanceModifier{
		Gender:        gender,
		ColorModifier: 1.0,
		SizeModifier:  1.0,
	}
}

// ApplyGenderModifier 应用性别修饰到外观
func ApplyGenderModifier(appearance Appearance, gender Gender) Appearance {
	modifier := GetGenderModifier(gender)

	// 修改体型
	if modifier.SizeModifier != 1.0 {
		// 体型修饰只影响显示，不改变基因值
		// 这里保持原值，实际渲染时应用修饰
	}

	return appearance
}
