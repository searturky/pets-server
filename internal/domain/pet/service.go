// Package pet 宠物领域
// 领域服务 - 处理跨实体的复杂业务逻辑
package pet

// DomainService 宠物领域服务
type DomainService struct {
	repo            Repository
	speciesRegistry *SpeciesRegistry
	fusionRegistry  *SpeciesFusionRegistry
	breedingService *BreedingService
}

// NewDomainService 创建领域服务
func NewDomainService(repo Repository, speciesRegistry *SpeciesRegistry, fusionRegistry *SpeciesFusionRegistry) *DomainService {
	return &DomainService{
		repo:            repo,
		speciesRegistry: speciesRegistry,
		fusionRegistry:  fusionRegistry,
		breedingService: NewBreedingService(speciesRegistry, fusionRegistry),
	}
}

// GetSpeciesRegistry 获取物种注册表
func (s *DomainService) GetSpeciesRegistry() *SpeciesRegistry {
	return s.speciesRegistry
}

// GetFusionRegistry 获取融合注册表
func (s *DomainService) GetFusionRegistry() *SpeciesFusionRegistry {
	return s.fusionRegistry
}

// GetBreedingService 获取繁衍服务
func (s *DomainService) GetBreedingService() *BreedingService {
	return s.breedingService
}

// --- 宠物创建 ---

// CreatePet 创建新宠物（指定物种）
func (s *DomainService) CreatePet(userID int64, name string, speciesID SpeciesID) (*Pet, error) {
	// 获取物种定义
	species, ok := s.speciesRegistry.Get(speciesID)
	if !ok {
		return nil, ErrSpeciesNotFound
	}

	// 创建宠物
	pet := NewPetWithSpecies(userID, name, speciesID, &species.GenderRule)

	// 解析物种特有外观
	if interpreter, ok := s.speciesRegistry.GetInterpreter(speciesID); ok {
		pet.SetSpecialAppearance(interpreter.InterpretSpecialFeatures(pet.Gene))
	}

	return pet, nil
}

// CreateRandomPet 创建随机物种的宠物
func (s *DomainService) CreateRandomPet(userID int64, name string) *Pet {
	// 获取可用物种列表
	availableSpecies := s.speciesRegistry.GetAvailableSpecies()
	if len(availableSpecies) == 0 {
		return NewPet(userID, name)
	}

	// 按稀有度加权随机选择
	totalWeight := 0
	for _, sp := range availableSpecies {
		// 稀有度越低，权重越高
		weight := 6 - sp.Rarity
		if weight < 1 {
			weight = 1
		}
		totalWeight += weight
	}

	roll := randomInt(totalWeight)
	accumulated := 0
	var selectedSpecies *Species
	for _, sp := range availableSpecies {
		weight := 6 - sp.Rarity
		if weight < 1 {
			weight = 1
		}
		accumulated += weight
		if roll < accumulated {
			selectedSpecies = sp
			break
		}
	}

	if selectedSpecies == nil {
		selectedSpecies = availableSpecies[0]
	}

	pet := NewPetWithSpecies(userID, name, selectedSpecies.ID, &selectedSpecies.GenderRule)

	// 解析物种特有外观
	if interpreter, ok := s.speciesRegistry.GetInterpreter(selectedSpecies.ID); ok {
		pet.SetSpecialAppearance(interpreter.InterpretSpecialFeatures(pet.Gene))
	}

	return pet
}

// --- 繁殖相关 ---

// BreedPets 繁殖两只宠物
func (s *DomainService) BreedPets(parent1, parent2 *Pet, childName string, ownerID int64) (*BreedingResult, error) {
	return s.breedingService.Breed(BreedingRequest{
		Parent1:   parent1,
		Parent2:   parent2,
		ChildName: childName,
		OwnerID:   ownerID,
	})
}

// SelfBreedPet 分裂繁殖
func (s *DomainService) SelfBreedPet(parent *Pet, childName string, ownerID int64) (*BreedingResult, error) {
	return s.breedingService.Breed(BreedingRequest{
		Parent1:   parent,
		Parent2:   nil,
		ChildName: childName,
		OwnerID:   ownerID,
	})
}

// CanBreedPair 检查两只宠物是否可以繁殖
func (s *DomainService) CanBreedPair(parent1, parent2 *Pet) error {
	return s.breedingService.CanBreedPair(parent1, parent2)
}

// CanSelfBreed 检查宠物是否可以分裂繁殖
func (s *DomainService) CanSelfBreed(pet *Pet) error {
	return s.breedingService.CanSelfBreed(pet)
}

// PredictOffspringSpecies 预测后代物种
func (s *DomainService) PredictOffspringSpecies(parent1, parent2 *Pet) []SpeciesProbability {
	return s.breedingService.PredictOffspringSpecies(parent1, parent2)
}

// --- 评分计算 ---

// CalculatePetScore 计算宠物综合评分
func (s *DomainService) CalculatePetScore(pet *Pet) int {
	score := 0

	// 基础分：等级
	score += pet.Level * 10

	// 稀有度加分：技能强度
	score += pet.Skill.Strength * 50

	// 物种稀有度加分
	if species, ok := s.speciesRegistry.Get(pet.SpeciesID); ok {
		score += species.Rarity * 100
	}

	// 状态加分
	score += (pet.Hunger + pet.Happiness + pet.Cleanliness) / 3

	// 进化阶段加分
	score += int(pet.Stage) * 100

	// 代数加分（越高代数越珍贵）
	score += pet.Generation * 20

	return score
}

// --- 物种查询 ---

// GetSpecies 获取物种信息
func (s *DomainService) GetSpecies(id SpeciesID) (*Species, bool) {
	return s.speciesRegistry.Get(id)
}

// GetAvailableSpecies 获取可用物种列表
func (s *DomainService) GetAvailableSpecies() []*Species {
	return s.speciesRegistry.GetAvailableSpecies()
}

// GetSpeciesByCategory 按分类获取物种
func (s *DomainService) GetSpeciesByCategory(category SpeciesCategory) []*Species {
	return s.speciesRegistry.GetByCategory(category)
}

// InterpretPetAppearance 解析宠物的物种特有外观
func (s *DomainService) InterpretPetAppearance(pet *Pet) SpecialAppearance {
	return s.speciesRegistry.InterpretGene(pet.SpeciesID, pet.Gene)
}
