// Package pet 宠物领域
// Personality 性格值对象 - 影响宠物行为和游戏机制
package pet

import "strings"

// Personality 性格值对象（不可变）
// 表示宠物的性格特质，影响游戏机制
type Personality struct {
	Activity  int // 活跃度 0-100: 影响快乐衰减速度和玩耍经验
	Appetite  int // 贪吃度 0-100: 影响饥饿衰减速度和喂食经验
	Social    int // 社交度 0-100: 影响拜访奖励
	Curiosity int // 好奇度 0-100: 影响随机事件触发概率
}

// NewPersonalityFromGene 从基因创建性格
func NewPersonalityFromGene(gene Gene) Personality {
	return Personality{
		Activity:  gene.ActivityTrait(),
		Appetite:  gene.AppetiteTrait(),
		Social:    gene.SocialTrait(),
		Curiosity: gene.CuriosityTrait(),
	}
}

// HungerDecayRate 饥饿衰减速率倍数
// 贪吃的宠物饿得更快
func (p Personality) HungerDecayRate() float64 {
	if p.Appetite > 70 {
		return 1.3
	} else if p.Appetite < 30 {
		return 0.8
	}
	return 1.0
}

// HappinessDecayRate 快乐衰减速率倍数
// 活跃的宠物需要更多互动
func (p Personality) HappinessDecayRate() float64 {
	if p.Activity > 70 {
		return 1.3
	} else if p.Activity < 30 {
		return 0.8
	}
	return 1.0
}

// FeedExpBonus 喂食经验加成倍数
func (p Personality) FeedExpBonus() float64 {
	if p.Appetite > 70 {
		return 1.2
	}
	return 1.0
}

// PlayExpBonus 玩耍经验加成倍数
func (p Personality) PlayExpBonus() float64 {
	if p.Activity > 70 {
		return 1.2
	}
	return 1.0
}

// VisitBonus 拜访奖励加成倍数
func (p Personality) VisitBonus() float64 {
	if p.Social > 70 {
		return 1.5
	} else if p.Social > 50 {
		return 1.2
	}
	return 1.0
}

// RandomEventChance 随机事件触发概率加成
func (p Personality) RandomEventChance() float64 {
	return 1.0 + float64(p.Curiosity)/200.0 // 0-50% 加成
}

// Describe 性格描述
func (p Personality) Describe() string {
	traits := []string{}

	if p.Activity > 70 {
		traits = append(traits, "活泼好动")
	} else if p.Activity < 30 {
		traits = append(traits, "安静沉稳")
	}

	if p.Appetite > 70 {
		traits = append(traits, "贪吃")
	} else if p.Appetite < 30 {
		traits = append(traits, "挑食")
	}

	if p.Social > 70 {
		traits = append(traits, "爱交朋友")
	} else if p.Social < 30 {
		traits = append(traits, "害羞")
	}

	if p.Curiosity > 70 {
		traits = append(traits, "好奇心强")
	} else if p.Curiosity < 30 {
		traits = append(traits, "谨慎")
	}

	if len(traits) == 0 {
		return "性格温和"
	}
	return strings.Join(traits, "、")
}

