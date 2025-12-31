// Package pet 宠物领域
// Appearance 外观值对象 - 由基因解码得到
package pet

import "fmt"

// Appearance 外观值对象（不可变）
// 表示宠物的外观属性，从基因解码而来
type Appearance struct {
	ColorPrimary   string // 主色 #RRGGBB
	ColorSecondary string // 副色 #RRGGBB
	PatternType    int    // 花纹类型 0-7
	BodyType       int    // 体型 0-3
	EarType        int    // 耳朵类型 0-5
	TailType       int    // 尾巴类型 0-7
	EyeType        int    // 眼睛类型 0-7
}

// NewAppearanceFromGene 从基因创建外观
func NewAppearanceFromGene(gene Gene) Appearance {
	return Appearance{
		ColorPrimary:   hslToHex(gene.PrimaryColorHue(), gene.PrimarySaturation(), 50),
		ColorSecondary: hslToHex(gene.SecondaryColorHue(), 50, 50),
		PatternType:    gene.PatternType(),
		BodyType:       gene.BodyType(),
		EarType:        gene.EarType(),
		TailType:       gene.TailType(),
		EyeType:        gene.EyeType(),
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
	return fmt.Sprintf("%s体型，%s花纹", a.BodyTypeName(), a.PatternTypeName())
}

