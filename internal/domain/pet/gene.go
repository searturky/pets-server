// Package pet 宠物领域
// Gene 基因值对象 - 控制宠物的随机性和独特性
package pet

import (
	"crypto/rand"
	"encoding/hex"
	"strconv"
)

// Gene 基因值对象（不可变）
// 20位十六进制字符串，决定宠物的外观、性格、技能等
// 结构：
//   - 位置 1-4:  外观A（主色、副色、体型）
//   - 位置 5-8:  外观B（花纹、耳朵、尾巴、眼睛）
//   - 位置 9-12: 性格（活跃度、贪吃度、社交度、好奇度）
//   - 位置 13-16: 技能（技能ID、技能强度）
//   - 位置 17-20: 隐藏（突变因子、进化倾向）
type Gene struct {
	code string // 20位十六进制
}

// NewGene 从字符串创建基因
func NewGene(code string) Gene {
	if len(code) != 20 {
		return GenerateGene() // 无效时生成新基因
	}
	return Gene{code: code}
}

// GenerateGene 生成随机基因
func GenerateGene() Gene {
	bytes := make([]byte, 10)
	rand.Read(bytes)
	return Gene{code: hex.EncodeToString(bytes)}
}

// String 返回基因码字符串
func (g Gene) String() string {
	return g.code
}

// --- 外观解析 ---

// PrimaryColorHue 主色色相 (0-360)
func (g Gene) PrimaryColorHue() int {
	return g.hexAt(0) * 24 // 0-15 映射到 0-360
}

// PrimarySaturation 主色饱和度 (0-100)
func (g Gene) PrimarySaturation() int {
	return int(float64(g.hexAt(1)) * 6.67)
}

// SecondaryColorHue 副色色相 (0-360)
func (g Gene) SecondaryColorHue() int {
	return g.hexAt(2) * 24
}

// BodyType 体型 (0-3: 小/中小/中/大)
func (g Gene) BodyType() int {
	return g.hexAt(3) % 4
}

// PatternType 花纹类型 (0-7)
func (g Gene) PatternType() int {
	return g.hexAt(4) % 8
}

// EarType 耳朵类型 (0-5)
func (g Gene) EarType() int {
	return g.hexAt(5) % 6
}

// TailType 尾巴类型 (0-7)
func (g Gene) TailType() int {
	return g.hexAt(6) % 8
}

// EyeType 眼睛类型 (0-7)
func (g Gene) EyeType() int {
	return g.hexAt(7) % 8
}

// --- 性格解析 ---

// ActivityTrait 活跃度 (0-100)
func (g Gene) ActivityTrait() int {
	return int(float64(g.hexAt(8)) * 6.67)
}

// AppetiteTrait 贪吃度 (0-100)
func (g Gene) AppetiteTrait() int {
	return int(float64(g.hexAt(9)) * 6.67)
}

// SocialTrait 社交度 (0-100)
func (g Gene) SocialTrait() int {
	return int(float64(g.hexAt(10)) * 6.67)
}

// CuriosityTrait 好奇度 (0-100)
func (g Gene) CuriosityTrait() int {
	return int(float64(g.hexAt(11)) * 6.67)
}

// --- 技能解析 ---

// SkillID 技能ID
func (g Gene) SkillID(totalSkills int) int {
	raw := g.hexAt(12)*16 + g.hexAt(13) // 0-255
	return raw % totalSkills
}

// SkillStrength 技能强度 (1-5星)
func (g Gene) SkillStrength() int {
	raw := g.hexAt(14)*16 + g.hexAt(15) // 0-255
	return (raw / 52) + 1               // 映射到 1-5
}

// --- 隐藏基因 ---

// MutationFactor 突变因子 (0-255)
func (g Gene) MutationFactor() int {
	return g.hexAt(16)*16 + g.hexAt(17)
}

// EvolutionTendency 进化倾向 (0-255)
func (g Gene) EvolutionTendency() int {
	return g.hexAt(18)*16 + g.hexAt(19)
}

// hexAt 获取指定位置的十六进制值 (0-15)
func (g Gene) hexAt(pos int) int {
	if pos >= len(g.code) {
		return 0
	}
	val, _ := strconv.ParseInt(string(g.code[pos]), 16, 64)
	return int(val)
}

// InheritFrom 基因遗传（如果将来支持繁殖）
func InheritFrom(parent1, parent2 Gene) Gene {
	child := make([]byte, 20)
	for i := 0; i < 20; i++ {
		// 50% 概率继承父母任一方
		randByte := make([]byte, 1)
		rand.Read(randByte)
		if randByte[0]%2 == 0 {
			child[i] = parent1.code[i]
		} else {
			child[i] = parent2.code[i]
		}
		// 5% 概率突变
		if randByte[0]%20 == 0 {
			hexChars := "0123456789abcdef"
			child[i] = hexChars[randByte[0]%16]
		}
	}
	return Gene{code: string(child)}
}

