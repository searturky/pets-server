// Package pet 宠物领域
// Appearance 外观值对象 - 由基因解码得到
package pet

import "fmt"

// Appearance 通用外观值对象（不可变）
// 表示所有物种共有的外观属性，从基因通用外观区域（0-7）解码而来
type Appearance struct {
	ColorPrimary   string // 主色 #RRGGBB
	ColorSecondary string // 副色 #RRGGBB
	PatternType    int    // 花纹类型 0-7
	PatternDensity int    // 花纹密度 0-15
	BodyType       int    // 体型 0-3
	EyeShape       int    // 眼睛形状 0-7
	EyeColor       int    // 眼睛颜色 0-15
}

// NewAppearanceFromGene 从基因创建通用外观
func NewAppearanceFromGene(gene Gene) Appearance {
	return Appearance{
		ColorPrimary:   hslToHex(gene.PrimaryColorHue(), gene.PrimarySaturation(), 50),
		ColorSecondary: hslToHex(gene.SecondaryColorHue(), 50, 50),
		PatternType:    gene.BasePattern(),
		PatternDensity: gene.PatternDensity(),
		BodyType:       gene.BodyType(),
		EyeShape:       gene.EyeShape(),
		EyeColor:       gene.EyeColor(),
	}
}

// BodyTypeName 体型名称
func (a Appearance) BodyTypeName() string {
	names := []string{"娇小", "小巧", "中等", "壮硕"}
	if a.BodyType >= 0 && a.BodyType < len(names) {
		return names[a.BodyType]
	}
	return "未知"
}

// PatternTypeName 花纹名称
func (a Appearance) PatternTypeName() string {
	names := []string{"纯色", "斑点", "条纹", "渐变", "双色", "三色", "星点", "云纹"}
	if a.PatternType >= 0 && a.PatternType < len(names) {
		return names[a.PatternType]
	}
	return "未知"
}

// EyeShapeName 眼睛形状名称
func (a Appearance) EyeShapeName() string {
	names := []string{"圆眼", "杏眼", "凤眼", "猫眼", "大眼", "小眼", "细长眼", "水汪汪"}
	if a.EyeShape >= 0 && a.EyeShape < len(names) {
		return names[a.EyeShape]
	}
	return "未知"
}

// hslToHex HSL转HEX颜色
func hslToHex(h, s, l int) string {
	// 简化实现：直接根据色相返回预设颜色
	// 实际项目中应该实现完整的HSL到RGB转换
	hueColors := []string{
		"#FF6B6B", "#FF8E53", "#FFD93D", "#6BCB77",
		"#4D96FF", "#9B59B6", "#E91E63", "#00BCD4",
		"#8BC34A", "#FF5722", "#673AB7", "#2196F3",
		"#FF9800", "#9C27B0", "#3F51B5", "#F44336",
	}
	idx := (h / 24) % len(hueColors)
	return hueColors[idx]
}

// String 外观描述
func (a Appearance) String() string {
	return fmt.Sprintf("%s体型，%s花纹，%s", a.BodyTypeName(), a.PatternTypeName(), a.EyeShapeName())
}

// SpecialAppearance 物种特有外观值对象
// 表示物种特有的外观属性，由物种解释器解读基因物种特征区域（8-15）
type SpecialAppearance struct {
	Parts []PartAppearance // 物种特有部位列表
}

// PartAppearance 部位外观
type PartAppearance struct {
	PartType PartType // 部位类型
	Style    string   // 样式名称
	Value    int      // 基因原始值
	Modifier int      // 修饰值
}

// PartType 部位类型
type PartType int

const (
	PartTypeNone   PartType = 0
	PartTypeEar    PartType = 1  // 耳朵
	PartTypeTail   PartType = 2  // 尾巴
	PartTypeFur    PartType = 3  // 毛纹
	PartTypeWing   PartType = 4  // 翅膀
	PartTypeBeak   PartType = 5  // 喙
	PartTypeCrest  PartType = 6  // 羽冠
	PartTypeFin    PartType = 7  // 鳍
	PartTypeScale  PartType = 8  // 鳞片
	PartTypeTailFin PartType = 9  // 尾鳍
	PartTypeShell  PartType = 10 // 壳
	PartTypeHorn   PartType = 11 // 角
	PartTypeArmor  PartType = 12 // 鳞甲
	PartTypeAura   PartType = 13 // 光环
	PartTypeClaw   PartType = 14 // 爪子
	PartTypeWhisker PartType = 15 // 胡须
)

// PartTypeName 部位类型名称
func (p PartType) Name() string {
	names := map[PartType]string{
		PartTypeNone:    "无",
		PartTypeEar:     "耳朵",
		PartTypeTail:    "尾巴",
		PartTypeFur:     "毛纹",
		PartTypeWing:    "翅膀",
		PartTypeBeak:    "喙",
		PartTypeCrest:   "羽冠",
		PartTypeFin:     "鳍",
		PartTypeScale:   "鳞片",
		PartTypeTailFin: "尾鳍",
		PartTypeShell:   "壳",
		PartTypeHorn:    "角",
		PartTypeArmor:   "鳞甲",
		PartTypeAura:    "光环",
		PartTypeClaw:    "爪子",
		PartTypeWhisker: "胡须",
	}
	if name, ok := names[p]; ok {
		return name
	}
	return "未知"
}

// NewSpecialAppearance 创建空的物种特有外观
func NewSpecialAppearance() SpecialAppearance {
	return SpecialAppearance{
		Parts: make([]PartAppearance, 0),
	}
}

// AddPart 添加部位外观
func (s *SpecialAppearance) AddPart(part PartAppearance) {
	s.Parts = append(s.Parts, part)
}

// GetPart 获取指定类型的部位外观
func (s SpecialAppearance) GetPart(partType PartType) (PartAppearance, bool) {
	for _, part := range s.Parts {
		if part.PartType == partType {
			return part, true
		}
	}
	return PartAppearance{}, false
}

// String 物种特有外观描述
func (s SpecialAppearance) String() string {
	if len(s.Parts) == 0 {
		return "无特殊外观"
	}
	result := ""
	for i, part := range s.Parts {
		if i > 0 {
			result += "，"
		}
		result += fmt.Sprintf("%s: %s", part.PartType.Name(), part.Style)
	}
	return result
}
