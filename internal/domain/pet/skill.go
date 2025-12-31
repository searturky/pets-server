// Package pet 宠物领域
// Skill 技能值对象 - 宠物的独特能力
package pet

// SkillType 技能类型
type SkillType int

const (
	SkillTypeNone       SkillType = 0
	SkillTypeLucky      SkillType = 1 // 幸运：增加金币获取
	SkillTypeCharming   SkillType = 2 // 魅力：增加拜访奖励
	SkillTypeEndurance  SkillType = 3 // 耐力：减缓状态衰减
	SkillTypeGluttony   SkillType = 4 // 大胃王：喂食效果增强
	SkillTypePlayful    SkillType = 5 // 爱玩：玩耍效果增强
	SkillTypeCleanlover SkillType = 6 // 爱干净：清洁效果增强
	SkillTypeFriendly   SkillType = 7 // 友善：社交奖励增强
	SkillTypeCurious    SkillType = 8 // 探索：随机事件更多
)

// Skill 技能值对象
type Skill struct {
	Type     SkillType // 技能类型
	Level    int       // 技能等级 1-5
	Strength int       // 技能强度 (来自基因)
}

// NewSkillFromGene 从基因创建技能
func NewSkillFromGene(gene Gene) Skill {
	totalSkills := 9 // 技能总数
	skillID := gene.SkillID(totalSkills)
	strength := gene.SkillStrength()

	return Skill{
		Type:     SkillType(skillID),
		Level:    1,
		Strength: strength,
	}
}

// Name 技能名称
func (s Skill) Name() string {
	names := map[SkillType]string{
		SkillTypeNone:       "无",
		SkillTypeLucky:      "幸运星",
		SkillTypeCharming:   "万人迷",
		SkillTypeEndurance:  "铁打的",
		SkillTypeGluttony:   "大胃王",
		SkillTypePlayful:    "玩乐达人",
		SkillTypeCleanlover: "洁癖",
		SkillTypeFriendly:   "社交达人",
		SkillTypeCurious:    "探险家",
	}
	if name, ok := names[s.Type]; ok {
		return name
	}
	return "未知"
}

// Description 技能描述
func (s Skill) Description() string {
	descs := map[SkillType]string{
		SkillTypeNone:       "没有特殊技能",
		SkillTypeLucky:      "获得金币时有概率额外获得",
		SkillTypeCharming:   "被拜访时双方获得更多奖励",
		SkillTypeEndurance:  "状态衰减速度减缓",
		SkillTypeGluttony:   "喂食恢复效果增强",
		SkillTypePlayful:    "玩耍恢复效果增强",
		SkillTypeCleanlover: "清洁恢复效果增强",
		SkillTypeFriendly:   "赠送礼物时获得额外好感",
		SkillTypeCurious:    "更容易触发随机事件",
	}
	if desc, ok := descs[s.Type]; ok {
		return desc
	}
	return "未知技能"
}

// EffectMultiplier 技能效果倍数
// 基于技能等级和强度计算
func (s Skill) EffectMultiplier() float64 {
	base := 1.0 + float64(s.Level)*0.1        // 等级加成
	strengthBonus := float64(s.Strength) * 0.02 // 强度加成
	return base + strengthBonus
}

// LevelUp 升级技能
func (s *Skill) LevelUp() bool {
	if s.Level >= 5 {
		return false
	}
	s.Level++
	return true
}

// Rarity 技能稀有度描述
func (s Skill) Rarity() string {
	if s.Strength >= 4 {
		return "传说"
	} else if s.Strength >= 3 {
		return "史诗"
	} else if s.Strength >= 2 {
		return "稀有"
	}
	return "普通"
}

