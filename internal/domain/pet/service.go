// Package pet 宠物领域
// 领域服务 - 处理跨实体的复杂业务逻辑
package pet

// TODO: 领域服务
// 当业务逻辑涉及多个实体或不属于单个实体时，放在领域服务中
// 例如：
// - 宠物繁殖（涉及两个宠物的基因遗传）
// - 宠物评分计算（综合多种属性）
// - 宠物对战（涉及两个宠物）

// DomainService 宠物领域服务
type DomainService struct {
	repo Repository
}

// NewDomainService 创建领域服务
func NewDomainService(repo Repository) *DomainService {
	return &DomainService{repo: repo}
}

// CalculatePetScore 计算宠物综合评分
// 这是一个不属于单个实体的业务逻辑示例
func (s *DomainService) CalculatePetScore(pet *Pet) int {
	score := 0

	// 基础分：等级
	score += pet.Level * 10

	// 稀有度加分：技能强度
	score += pet.Skill.Strength * 50

	// 状态加分
	score += (pet.Hunger + pet.Happiness + pet.Cleanliness) / 3

	// 进化阶段加分
	score += int(pet.Stage) * 100

	return score
}

// BreedPets 繁殖两只宠物（示例，暂未实现完整逻辑）
// 返回新宠物的基因
func (s *DomainService) BreedPets(pet1, pet2 *Pet) Gene {
	return InheritFrom(pet1.Gene, pet2.Gene)
}

