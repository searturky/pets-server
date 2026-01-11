// Package pet 宠物领域
// Personality 性格值对象 - 影响宠物行为和游戏机制
package pet

import "strings"

// Personality 性格值对象（不可变）
// 表示宠物的性格特质，影响游戏机制
// 从基因性格区域（16-23）解码而来
type Personality struct {
	Activity     int // 活跃度 0-100: 影响快乐衰减速度和玩耍经验
	Appetite     int // 贪吃度 0-100: 影响饥饿衰减速度和喂食经验
	Social       int // 社交度 0-100: 影响拜访奖励
	Curiosity    int // 好奇度 0-100: 影响随机事件触发概率
	Temper       int // 脾气 0-100: 影响互动反应
	Loyalty      int // 忠诚度 0-100: 影响与主人的互动
	Intelligence int // 智力 0-100: 影响学习速度
	Playfulness  int // 玩乐度 0-100: 影响玩耍效果
}

// NewPersonalityFromGene 从基因创建性格
func NewPersonalityFromGene(gene Gene) Personality {
	return Personality{
		Activity:     gene.ActivityTrait(),
		Appetite:     gene.AppetiteTrait(),
		Social:       gene.SocialTrait(),
		Curiosity:    gene.CuriosityTrait(),
		Temper:       gene.TemperTrait(),
		Loyalty:      gene.LoyaltyTrait(),
		Intelligence: gene.IntelligenceTrait(),
		Playfulness:  gene.PlayfulnessTrait(),
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
	bonus := 1.0
	if p.Activity > 70 {
		bonus += 0.1
	}
	if p.Playfulness > 70 {
		bonus += 0.1
	}
	return bonus
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

// LearningSpeed 学习速度倍数
func (p Personality) LearningSpeed() float64 {
	return 1.0 + float64(p.Intelligence)/200.0 // 0-50% 加成
}

// LoyaltyBonus 忠诚度奖励加成
func (p Personality) LoyaltyBonus() float64 {
	if p.Loyalty > 70 {
		return 1.3
	} else if p.Loyalty > 50 {
		return 1.15
	}
	return 1.0
}

// TemperReaction 脾气反应类型
func (p Personality) TemperReaction() string {
	if p.Temper > 80 {
		return "暴躁"
	} else if p.Temper > 60 {
		return "易怒"
	} else if p.Temper > 40 {
		return "温和"
	} else if p.Temper > 20 {
		return "温顺"
	}
	return "冷静"
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

	if p.Intelligence > 70 {
		traits = append(traits, "聪明")
	} else if p.Intelligence < 30 {
		traits = append(traits, "呆萌")
	}

	if p.Playfulness > 70 {
		traits = append(traits, "爱玩")
	}

	if p.Loyalty > 70 {
		traits = append(traits, "忠诚")
	}

	if len(traits) == 0 {
		return "性格温和"
	}
	return strings.Join(traits, "、")
}
