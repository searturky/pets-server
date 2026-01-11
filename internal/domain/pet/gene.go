// Package pet 宠物领域
// Gene 基因值对象 - 控制宠物的随机性和独特性
package pet

import (
	"crypto/rand"
	"encoding/hex"
	"strconv"
)

// GeneLength 基因码长度（40位十六进制 = 160 bits）
const GeneLength = 40

// 基因片段位置常量
const (
	// 通用外观 (0-7)
	GenePosColorPrimary    = 0 // 主色色相
	GenePosColorSaturation = 1 // 主色饱和度
	GenePosColorSecondary  = 2 // 副色色相
	GenePosBodySize        = 3 // 体型大小
	GenePosEyeShape        = 4 // 眼睛形状
	GenePosEyeColor        = 5 // 眼睛颜色
	GenePosBasePattern     = 6 // 基础花纹
	GenePosPatternDensity  = 7 // 花纹密度

	// 物种特征 (8-15) - 由物种解释器解读
	GenePosSpecialA    = 8  // 特征A (如:耳朵/翅膀/鳍)
	GenePosSpecialB    = 9  // 特征B (如:尾巴/喙/鳞)
	GenePosSpecialC    = 10 // 特征C (如:毛纹/羽冠/甲壳)
	GenePosSpecialD    = 11 // 特征D (预留/复合特征)
	GenePosSpecialModA = 12 // 特征A修饰
	GenePosSpecialModB = 13 // 特征B修饰
	GenePosSpecialModC = 14 // 特征C修饰
	GenePosSpecialModD = 15 // 特征D修饰

	// 性格 (16-23)
	GenePosActivity     = 16 // 活跃度
	GenePosAppetite     = 17 // 贪吃度
	GenePosSocial       = 18 // 社交度
	GenePosCuriosity    = 19 // 好奇度
	GenePosTemper       = 20 // 脾气
	GenePosLoyalty      = 21 // 忠诚度
	GenePosIntelligence = 22 // 智力
	GenePosPlayfulness  = 23 // 玩乐度

	// 技能/能力 (24-31)
	GenePosSkillPrimary   = 24 // 主技能ID
	GenePosSkillStrength  = 25 // 主技能强度
	GenePosSkillSecondary = 26 // 副技能ID
	GenePosAbilityA       = 27 // 特殊能力A
	GenePosAbilityB       = 28 // 特殊能力B
	GenePosGrowthRate     = 29 // 成长速率
	GenePosLifespan       = 30 // 寿命倾向
	GenePosResistance     = 31 // 抗性

	// 遗传隐藏 (32-39)
	GenePosMutation      = 32 // 突变因子
	GenePosEvolution     = 33 // 进化倾向
	GenePosHiddenSpecies = 34 // 隐藏物种触发
	GenePosHiddenTrait   = 35 // 隐藏特质
	GenePosBreedBonus    = 36 // 繁殖加成
	GenePosRecessiveA    = 37 // 隐性基因A
	GenePosRecessiveB    = 38 // 隐性基因B
	GenePosRecessiveC    = 39 // 隐性基因C
)

// Gene 基因值对象（不可变）
// 40位十六进制字符串，决定宠物的外观、性格、技能等
// 结构：
//   - 位置 0-7:   通用外观（主色、副色、体型、眼睛、花纹）
//   - 位置 8-15:  物种特征（由物种解释器解读）
//   - 位置 16-23: 性格（活跃度、贪吃度、社交度等）
//   - 位置 24-31: 技能/能力（技能ID、强度、特殊能力）
//   - 位置 32-39: 遗传隐藏（突变、进化、隐藏物种、隐性基因）
type Gene struct {
	code string // 40位十六进制
}

// NewGene 从字符串创建基因
func NewGene(code string) Gene {
	if len(code) != GeneLength {
		return GenerateGene() // 无效时生成新基因
	}
	return Gene{code: code}
}

// GenerateGene 生成随机基因
func GenerateGene() Gene {
	bytes := make([]byte, GeneLength/2) // 20 bytes = 40 hex chars
	rand.Read(bytes)
	return Gene{code: hex.EncodeToString(bytes)}
}

// String 返回基因码字符串
func (g Gene) String() string {
	return g.code
}

// HexAt 获取指定位置的十六进制值 (0-15)
func (g Gene) HexAt(pos int) int {
	if pos >= len(g.code) || pos < 0 {
		return 0
	}
	val, _ := strconv.ParseInt(string(g.code[pos]), 16, 64)
	return int(val)
}

// HexPairAt 获取两个连续位置的组合值 (0-255)
func (g Gene) HexPairAt(pos int) int {
	return g.HexAt(pos)*16 + g.HexAt(pos+1)
}

// --- 通用外观解析 (位置 0-7) ---

// PrimaryColorHue 主色色相 (0-360)
func (g Gene) PrimaryColorHue() int {
	return g.HexAt(GenePosColorPrimary) * 24 // 0-15 映射到 0-360
}

// PrimarySaturation 主色饱和度 (0-100)
func (g Gene) PrimarySaturation() int {
	return int(float64(g.HexAt(GenePosColorSaturation)) * 6.67)
}

// SecondaryColorHue 副色色相 (0-360)
func (g Gene) SecondaryColorHue() int {
	return g.HexAt(GenePosColorSecondary) * 24
}

// BodySize 体型大小 (0-15)
func (g Gene) BodySize() int {
	return g.HexAt(GenePosBodySize)
}

// BodyType 体型类型 (0-3: 娇小/小巧/中等/壮硕)
func (g Gene) BodyType() int {
	return g.HexAt(GenePosBodySize) % 4
}

// EyeShape 眼睛形状 (0-7)
func (g Gene) EyeShape() int {
	return g.HexAt(GenePosEyeShape) % 8
}

// EyeColor 眼睛颜色 (0-15)
func (g Gene) EyeColor() int {
	return g.HexAt(GenePosEyeColor)
}

// BasePattern 基础花纹类型 (0-7)
func (g Gene) BasePattern() int {
	return g.HexAt(GenePosBasePattern) % 8
}

// PatternDensity 花纹密度 (0-15)
func (g Gene) PatternDensity() int {
	return g.HexAt(GenePosPatternDensity)
}

// --- 物种特征 (位置 8-15) - 原始值，由物种解释器解读 ---

// SpecialA 特征A原始值 (0-15)
func (g Gene) SpecialA() int {
	return g.HexAt(GenePosSpecialA)
}

// SpecialB 特征B原始值 (0-15)
func (g Gene) SpecialB() int {
	return g.HexAt(GenePosSpecialB)
}

// SpecialC 特征C原始值 (0-15)
func (g Gene) SpecialC() int {
	return g.HexAt(GenePosSpecialC)
}

// SpecialD 特征D原始值 (0-15)
func (g Gene) SpecialD() int {
	return g.HexAt(GenePosSpecialD)
}

// SpecialModA 特征A修饰值 (0-15)
func (g Gene) SpecialModA() int {
	return g.HexAt(GenePosSpecialModA)
}

// SpecialModB 特征B修饰值 (0-15)
func (g Gene) SpecialModB() int {
	return g.HexAt(GenePosSpecialModB)
}

// SpecialModC 特征C修饰值 (0-15)
func (g Gene) SpecialModC() int {
	return g.HexAt(GenePosSpecialModC)
}

// SpecialModD 特征D修饰值 (0-15)
func (g Gene) SpecialModD() int {
	return g.HexAt(GenePosSpecialModD)
}

// SpecialWithMod 获取特征值与修饰值的组合 (0-255)
func (g Gene) SpecialWithMod(featurePos, modPos int) int {
	return g.HexAt(featurePos)*16 + g.HexAt(modPos)
}

// --- 性格解析 (位置 16-23) ---

// ActivityTrait 活跃度 (0-100)
func (g Gene) ActivityTrait() int {
	return int(float64(g.HexAt(GenePosActivity)) * 6.67)
}

// AppetiteTrait 贪吃度 (0-100)
func (g Gene) AppetiteTrait() int {
	return int(float64(g.HexAt(GenePosAppetite)) * 6.67)
}

// SocialTrait 社交度 (0-100)
func (g Gene) SocialTrait() int {
	return int(float64(g.HexAt(GenePosSocial)) * 6.67)
}

// CuriosityTrait 好奇度 (0-100)
func (g Gene) CuriosityTrait() int {
	return int(float64(g.HexAt(GenePosCuriosity)) * 6.67)
}

// TemperTrait 脾气 (0-100)
func (g Gene) TemperTrait() int {
	return int(float64(g.HexAt(GenePosTemper)) * 6.67)
}

// LoyaltyTrait 忠诚度 (0-100)
func (g Gene) LoyaltyTrait() int {
	return int(float64(g.HexAt(GenePosLoyalty)) * 6.67)
}

// IntelligenceTrait 智力 (0-100)
func (g Gene) IntelligenceTrait() int {
	return int(float64(g.HexAt(GenePosIntelligence)) * 6.67)
}

// PlayfulnessTrait 玩乐度 (0-100)
func (g Gene) PlayfulnessTrait() int {
	return int(float64(g.HexAt(GenePosPlayfulness)) * 6.67)
}

// --- 技能/能力解析 (位置 24-31) ---

// SkillPrimaryID 主技能ID (0-255)
func (g Gene) SkillPrimaryID() int {
	return g.HexPairAt(GenePosSkillPrimary)
}

// SkillID 技能ID（兼容旧接口）
func (g Gene) SkillID(totalSkills int) int {
	return g.SkillPrimaryID() % totalSkills
}

// SkillStrength 技能强度 (1-5星)
func (g Gene) SkillStrength() int {
	raw := g.HexAt(GenePosSkillStrength)
	return (raw / 3) + 1 // 映射到 1-5
}

// SkillSecondaryID 副技能ID (0-15)
func (g Gene) SkillSecondaryID() int {
	return g.HexAt(GenePosSkillSecondary)
}

// AbilityA 特殊能力A (0-15)
func (g Gene) AbilityA() int {
	return g.HexAt(GenePosAbilityA)
}

// AbilityB 特殊能力B (0-15)
func (g Gene) AbilityB() int {
	return g.HexAt(GenePosAbilityB)
}

// GrowthRate 成长速率倍数
func (g Gene) GrowthRate() float64 {
	raw := g.HexAt(GenePosGrowthRate)
	return 0.8 + float64(raw)*0.025 // 0.8 - 1.175
}

// LifespanModifier 寿命修正
func (g Gene) LifespanModifier() float64 {
	raw := g.HexAt(GenePosLifespan)
	return 0.9 + float64(raw)*0.02 // 0.9 - 1.2
}

// Resistance 抗性 (0-15)
func (g Gene) Resistance() int {
	return g.HexAt(GenePosResistance)
}

// --- 遗传隐藏解析 (位置 32-39) ---

// MutationFactor 突变因子 (0-255)
func (g Gene) MutationFactor() int {
	return g.HexPairAt(GenePosMutation)
}

// EvolutionTendency 进化倾向 (0-255)
func (g Gene) EvolutionTendency() int {
	return g.HexPairAt(GenePosEvolution)
}

// HiddenSpeciesTrigger 隐藏物种触发值 (0-255)
func (g Gene) HiddenSpeciesTrigger() int {
	return g.HexPairAt(GenePosHiddenSpecies)
}

// HiddenTrait 隐藏特质 (0-15)
func (g Gene) HiddenTrait() int {
	return g.HexAt(GenePosHiddenTrait)
}

// BreedBonus 繁殖加成 (0-15)
func (g Gene) BreedBonus() int {
	return g.HexAt(GenePosBreedBonus)
}

// RecessiveA 隐性基因A (0-15)
func (g Gene) RecessiveA() int {
	return g.HexAt(GenePosRecessiveA)
}

// RecessiveB 隐性基因B (0-15)
func (g Gene) RecessiveB() int {
	return g.HexAt(GenePosRecessiveB)
}

// RecessiveC 隐性基因C (0-15)
func (g Gene) RecessiveC() int {
	return g.HexAt(GenePosRecessiveC)
}

// --- 基因遗传 ---

// hexChar 将数值转换为十六进制字符
func hexChar(val int) byte {
	const hexChars = "0123456789abcdef"
	return hexChars[val%16]
}

// randomByte 生成随机字节
func randomByte() byte {
	b := make([]byte, 1)
	rand.Read(b)
	return b[0]
}

// randomInt 生成 0 到 max-1 的随机整数
func randomInt(max int) int {
	if max <= 0 {
		return 0
	}
	return int(randomByte()) % max
}

// InheritFrom 基因遗传（有性繁殖）
// 按位从父母双方遗传，有突变和混合概率
func InheritFrom(parent1, parent2 Gene) Gene {
	child := make([]byte, GeneLength)

	for i := 0; i < GeneLength; i++ {
		p1Val := parent1.HexAt(i)
		p2Val := parent2.HexAt(i)

		// 决定继承方式
		inheritMode := randomInt(100)

		var childVal int
		switch {
		case inheritMode < 45:
			// 45% 继承父方
			childVal = p1Val
		case inheritMode < 90:
			// 45% 继承母方
			childVal = p2Val
		case inheritMode < 97:
			// 7% 混合（取平均值）
			childVal = (p1Val + p2Val) / 2
		default:
			// 3% 突变（完全随机）
			childVal = randomInt(16)
		}

		child[i] = hexChar(childVal)
	}

	// 处理隐性基因的特殊遗传
	processRecessiveGenes(parent1, parent2, child)

	return Gene{code: string(child)}
}

// processRecessiveGenes 处理隐性基因的特殊遗传规则
// 当父母双方的隐性基因相同或接近时，有更高概率表达
func processRecessiveGenes(p1, p2 Gene, child []byte) {
	recessivePositions := []int{GenePosRecessiveA, GenePosRecessiveB, GenePosRecessiveC}

	for _, pos := range recessivePositions {
		p1Val := p1.HexAt(pos)
		p2Val := p2.HexAt(pos)

		if p1Val == p2Val {
			// 双亲相同 -> 100% 继承
			child[pos] = hexChar(p1Val)
		} else if abs(p1Val-p2Val) <= 2 {
			// 接近 -> 75% 继承较强者
			if randomInt(4) != 0 {
				child[pos] = hexChar(maxInt(p1Val, p2Val))
			}
		}
	}
}

// SelfReplicate 自我复制基因（无性繁殖/分裂）
// mutationRate: 突变率 (0.0 - 1.0)
func SelfReplicate(parent Gene, mutationRate float64) Gene {
	child := make([]byte, GeneLength)
	copy(child, parent.code)

	mutationThreshold := int(mutationRate * 100)

	for i := 0; i < GeneLength; i++ {
		if randomInt(100) < mutationThreshold {
			child[i] = hexChar(randomInt(16))
		}
	}

	return Gene{code: string(child)}
}

// abs 绝对值
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// maxInt 返回较大值
func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// minInt 返回较小值
func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
