// Package pet 宠物领域
// Breeding 繁衍系统 - 处理宠物繁殖的领域逻辑
package pet

import (
	"time"
)

// BreedingService 繁衍领域服务
type BreedingService struct {
	speciesRegistry *SpeciesRegistry
	fusionRegistry  *SpeciesFusionRegistry
}

// NewBreedingService 创建繁衍服务
func NewBreedingService(speciesRegistry *SpeciesRegistry, fusionRegistry *SpeciesFusionRegistry) *BreedingService {
	return &BreedingService{
		speciesRegistry: speciesRegistry,
		fusionRegistry:  fusionRegistry,
	}
}

// BreedingRequest 繁殖请求
type BreedingRequest struct {
	Parent1   *Pet   // 父方
	Parent2   *Pet   // 母方 (分裂繁殖时为nil)
	ChildName string // 子代名称
	OwnerID   int64  // 子代所有者ID
}

// BreedingResult 繁殖结果
type BreedingResult struct {
	Child       *Pet      // 子代宠物
	SpeciesID   SpeciesID // 子代物种
	Gender      Gender    // 子代性别
	IsHidden    bool      // 是否触发隐藏物种
	FusionFrom  []SpeciesID // 融合来源物种
}

// Breed 执行繁殖
func (s *BreedingService) Breed(req BreedingRequest) (*BreedingResult, error) {
	// 检查是否为分裂繁殖
	if req.Parent2 == nil {
		return s.selfBreed(req)
	}
	return s.sexualBreed(req)
}

// selfBreed 分裂繁殖（无性繁殖）
func (s *BreedingService) selfBreed(req BreedingRequest) (*BreedingResult, error) {
	parent := req.Parent1

	// 获取物种定义
	species, ok := s.speciesRegistry.Get(parent.SpeciesID)
	if !ok {
		return nil, ErrSpeciesNotFound
	}

	// 检查是否支持分裂繁殖
	if !species.GenderRule.CanSelfBreed {
		return nil, ErrCannotSelfBreed
	}

	// 检查繁殖条件
	if err := parent.CanBreed(species.BreedRules); err != nil {
		return nil, err
	}

	// 自我复制基因（15%突变率）
	childGene := SelfReplicate(parent.Gene, 0.15)

	// 确定性别（分裂繁殖继承相同性别/无性别）
	childGender := parent.Gender

	// 创建子代
	child := NewPetFromBreeding(
		req.OwnerID,
		req.ChildName,
		parent.SpeciesID,
		childGene,
		childGender,
		parent.ID,
		0, // 无第二亲本
		parent.Generation+1,
	)

	// 解析物种特有外观
	if interpreter, ok := s.speciesRegistry.GetInterpreter(parent.SpeciesID); ok {
		child.SetSpecialAppearance(interpreter.InterpretSpecialFeatures(childGene))
	}

	// 标记父方已繁殖
	parent.MarkBred()

	return &BreedingResult{
		Child:     child,
		SpeciesID: parent.SpeciesID,
		Gender:    childGender,
		IsHidden:  false,
	}, nil
}

// sexualBreed 有性繁殖
func (s *BreedingService) sexualBreed(req BreedingRequest) (*BreedingResult, error) {
	parent1 := req.Parent1
	parent2 := req.Parent2

	// 获取物种定义
	species1, ok := s.speciesRegistry.Get(parent1.SpeciesID)
	if !ok {
		return nil, ErrSpeciesNotFound
	}
	species2, ok := s.speciesRegistry.Get(parent2.SpeciesID)
	if !ok {
		return nil, ErrSpeciesNotFound
	}

	// 检查性别兼容性
	if err := parent1.CanBreedWith(parent2); err != nil {
		return nil, err
	}

	// 检查繁殖条件
	if err := parent1.CanBreed(species1.BreedRules); err != nil {
		return nil, err
	}
	if err := parent2.CanBreed(species2.BreedRules); err != nil {
		return nil, err
	}

	// 基因遗传
	childGene := InheritFrom(parent1.Gene, parent2.Gene)

	// 确定子代物种
	childSpeciesID, isHidden := s.determineChildSpecies(parent1, parent2, childGene)

	// 获取子代物种定义
	childSpecies, ok := s.speciesRegistry.Get(childSpeciesID)
	if !ok {
		// 如果隐藏物种不存在，回退到父方物种
		childSpeciesID = parent1.SpeciesID
		childSpecies = species1
	}

	// 确定子代性别
	childGender := DetermineChildGender(parent1, parent2, childGene, childSpecies.GenderRule)

	// 创建子代
	generation := maxInt(parent1.Generation, parent2.Generation) + 1
	child := NewPetFromBreeding(
		req.OwnerID,
		req.ChildName,
		childSpeciesID,
		childGene,
		childGender,
		parent1.ID,
		parent2.ID,
		generation,
	)

	// 解析物种特有外观
	if interpreter, ok := s.speciesRegistry.GetInterpreter(childSpeciesID); ok {
		child.SetSpecialAppearance(interpreter.InterpretSpecialFeatures(childGene))
	}

	// 标记双亲已繁殖
	parent1.MarkBred()
	parent2.MarkBred()

	// 构建结果
	result := &BreedingResult{
		Child:     child,
		SpeciesID: childSpeciesID,
		Gender:    childGender,
		IsHidden:  isHidden,
	}

	if isHidden {
		result.FusionFrom = []SpeciesID{parent1.SpeciesID, parent2.SpeciesID}
	}

	return result, nil
}

// determineChildSpecies 确定子代物种
func (s *BreedingService) determineChildSpecies(parent1, parent2 *Pet, childGene Gene) (SpeciesID, bool) {
	// 同物种繁殖 -> 100% 相同物种
	if parent1.SpeciesID == parent2.SpeciesID {
		return parent1.SpeciesID, false
	}

	// 跨物种繁殖
	roll := randomInt(100)

	switch {
	case roll < 45:
		// 45% 继承父方物种
		return parent1.SpeciesID, false
	case roll < 90:
		// 45% 继承母方物种
		return parent2.SpeciesID, false
	default:
		// 10% 触发隐藏物种检查
		if hiddenSpecies, triggered := s.fusionRegistry.CheckFusionTrigger(
			parent1.SpeciesID,
			parent2.SpeciesID,
			childGene,
		); triggered {
			return hiddenSpecies, true
		}
		// 未触发则随机选择父母物种
		if randomInt(2) == 0 {
			return parent1.SpeciesID, false
		}
		return parent2.SpeciesID, false
	}
}

// CanBreedPair 检查两只宠物是否可以繁殖
func (s *BreedingService) CanBreedPair(parent1, parent2 *Pet) error {
	// 获取物种定义
	species1, ok := s.speciesRegistry.Get(parent1.SpeciesID)
	if !ok {
		return ErrSpeciesNotFound
	}
	species2, ok := s.speciesRegistry.Get(parent2.SpeciesID)
	if !ok {
		return ErrSpeciesNotFound
	}

	// 检查性别兼容性
	if err := parent1.CanBreedWith(parent2); err != nil {
		return err
	}

	// 检查各自繁殖条件
	if err := parent1.CanBreed(species1.BreedRules); err != nil {
		return err
	}
	if err := parent2.CanBreed(species2.BreedRules); err != nil {
		return err
	}

	return nil
}

// CanSelfBreed 检查宠物是否可以分裂繁殖
func (s *BreedingService) CanSelfBreed(pet *Pet) error {
	species, ok := s.speciesRegistry.Get(pet.SpeciesID)
	if !ok {
		return ErrSpeciesNotFound
	}

	if !species.GenderRule.CanSelfBreed {
		return ErrCannotSelfBreed
	}

	return pet.CanBreed(species.BreedRules)
}

// GetBreedingCooldown 获取繁殖冷却剩余时间
func (s *BreedingService) GetBreedingCooldown(pet *Pet) time.Duration {
	if pet.LastBreedAt == nil {
		return 0
	}

	species, ok := s.speciesRegistry.Get(pet.SpeciesID)
	if !ok {
		return 0
	}

	var cooldownHours int
	if pet.Gender == GenderNone {
		cooldownHours = species.BreedRules.SelfBreedCooldownHours
	} else {
		cooldownHours = species.BreedRules.CooldownHours
	}

	cooldown := time.Duration(cooldownHours) * time.Hour
	elapsed := time.Since(*pet.LastBreedAt)

	if elapsed >= cooldown {
		return 0
	}
	return cooldown - elapsed
}

// PredictOffspringSpecies 预测后代可能的物种
func (s *BreedingService) PredictOffspringSpecies(parent1, parent2 *Pet) []SpeciesProbability {
	var result []SpeciesProbability

	// 同物种
	if parent1.SpeciesID == parent2.SpeciesID {
		return []SpeciesProbability{
			{SpeciesID: parent1.SpeciesID, Probability: 100},
		}
	}

	// 跨物种
	result = append(result, SpeciesProbability{
		SpeciesID:   parent1.SpeciesID,
		Probability: 45,
	})
	result = append(result, SpeciesProbability{
		SpeciesID:   parent2.SpeciesID,
		Probability: 45,
	})

	// 检查是否有隐藏物种
	if config, ok := s.fusionRegistry.GetFusion(parent1.SpeciesID, parent2.SpeciesID); ok {
		// 隐藏物种概率约为 10% * (触发概率)
		// 触发概率取决于基因阈值，这里简化为固定值
		hiddenProb := 10 * (255 - config.TriggerThreshold) / 255
		if hiddenProb < 1 {
			hiddenProb = 1
		}

		result = append(result, SpeciesProbability{
			SpeciesID:   config.ResultSpecies,
			Probability: hiddenProb,
			IsHidden:    true,
		})

		// 调整其他概率
		remaining := 100 - hiddenProb
		for i := range result[:2] {
			result[i].Probability = remaining * result[i].Probability / 90
		}
	}

	return result
}

// SpeciesProbability 物种概率
type SpeciesProbability struct {
	SpeciesID   SpeciesID
	Probability int  // 百分比
	IsHidden    bool // 是否为隐藏物种
}
