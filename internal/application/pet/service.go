// Package pet 宠物应用服务
// 编排宠物相关的业务用例
// 这是完整的读写示例
package pet

import (
	"context"
	"errors"
	"strconv"
	"time"

	"pets-server/internal/domain/item"
	"pets-server/internal/domain/pet"
	"pets-server/internal/domain/shared"
)

// Service 宠物应用服务
type Service struct {
	petRepo      pet.Repository
	itemRepo     item.Repository
	petDomainSvc *pet.DomainService // 领域服务
	uow          shared.UnitOfWork
	publisher    shared.EventPublisher
	cache        CacheService // 缓存服务接口
}

// CacheService 缓存服务接口（在应用层定义，基础设施层实现）
type CacheService interface {
	GetPetDetail(ctx context.Context, userID int64) (*PetDetailDTO, error)
	SetPetDetail(ctx context.Context, userID int64, pet *PetDetailDTO, ttl time.Duration) error
	DeletePetDetail(ctx context.Context, userID int64) error
}

// NewService 创建宠物应用服务
func NewService(
	petRepo pet.Repository,
	itemRepo item.Repository,
	petDomainSvc *pet.DomainService,
	uow shared.UnitOfWork,
	publisher shared.EventPublisher,
	cache CacheService,
) *Service {
	return &Service{
		petRepo:      petRepo,
		itemRepo:     itemRepo,
		petDomainSvc: petDomainSvc,
		uow:          uow,
		publisher:    publisher,
		cache:        cache,
	}
}

// ============================================================
// 读操作示例：获取宠物详情 (GET /api/pet)
// 调用链路：
//   Handler.GetMyPet()
//     → Redis.Get() 尝试读缓存
//     → 缓存未命中
//       → AppService.GetPetDetail()
//         → PetRepo.FindByUserID() 查询数据库
//         → 组装 DTO
//       → Redis.Set() 写入缓存
//     ← 返回 DTO
// ============================================================

// GetPetDetail 获取用户的宠物详情
func (s *Service) GetPetDetail(ctx context.Context, userID int64) (*PetDetailDTO, error) {
	// 1. 尝试从缓存获取
	if s.cache != nil {
		cached, err := s.cache.GetPetDetail(ctx, userID)
		if err == nil && cached != nil {
			return cached, nil
		}
		// 缓存未命中或出错，继续查询数据库
	}

	// 2. 从数据库获取宠物实体
	p, err := s.petRepo.FindByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, pet.ErrPetNotFound) {
			return nil, ErrPetNotFound
		}
		return nil, err
	}

	// 3. 将领域实体转换为 DTO
	dto := s.toPetDetailDTO(p)

	// 4. 写入缓存（异步或同步，这里用同步简化示例）
	if s.cache != nil {
		_ = s.cache.SetPetDetail(ctx, userID, dto, 5*time.Minute)
	}

	return dto, nil
}

// ============================================================
// 写操作示例：喂食宠物 (POST /api/pet/feed)
// 调用链路：
//   Handler.Feed()
//     → AppService.FeedPet()
//       → UoW.Do() 开启事务
//         → PetRepo.FindByUserID() 获取宠物
//         → ItemRepo.FindByUserAndItem() 获取食物
//         → Item.Consume() 领域逻辑：扣道具
//         → Pet.Feed() 领域逻辑：喂食
//         → ItemRepo.Save() 保存道具
//         → PetRepo.Save() 保存宠物
//       → 事务提交
//       → Cache.Delete() 清除缓存
//       → EventPublisher.Publish(PetFedEvent) 发布事件
//     ← 返回结果
// ============================================================

// FeedPet 喂食宠物
func (s *Service) FeedPet(ctx context.Context, userID int64, req FeedPetRequest) (*FeedPetResponse, error) {
	var response *FeedPetResponse
	var events []any

	// 1. 在事务中执行所有数据库操作
	err := s.uow.Do(ctx, func(txCtx context.Context) error {
		// 1.1 获取宠物实体
		p, err := s.petRepo.FindByUserID(txCtx, userID)
		if err != nil {
			return err
		}

		// 1.2 获取道具定义，确定食物类型
		itemDef, err := s.itemRepo.GetDefinition(txCtx, req.FoodItemID)
		if err != nil {
			return err
		}
		if itemDef.Type != item.ItemTypeFood {
			return ErrInvalidFoodItem
		}

		// 1.3 获取用户的该道具
		userItem, err := s.itemRepo.FindByUserAndItem(txCtx, userID, req.FoodItemID)
		if err != nil {
			return err
		}

		// 1.4 领域逻辑：消耗道具
		if err := userItem.Consume(1); err != nil {
			return err
		}

		// 1.5 记录喂食前的等级
		oldLevel := p.Level

		// 1.6 领域逻辑：喂食宠物
		foodType := pet.FoodType(itemDef.EffectValue)
		if err := p.Feed(foodType); err != nil {
			return err
		}

		// 1.7 保存道具变更
		if err := s.itemRepo.Save(txCtx, userItem); err != nil {
			return err
		}

		// 1.8 保存宠物变更
		if err := s.petRepo.Save(txCtx, p); err != nil {
			return err
		}

		// 1.9 收集领域事件
		events = p.Events()

		// 1.10 构建响应
		response = &FeedPetResponse{
			Hunger:    p.Hunger,
			ExpGained: 10, // 基础经验
			LevelUp:   p.Level > oldLevel,
			NewLevel:  p.Level,
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// 2. 事务成功后，清除缓存
	if s.cache != nil {
		_ = s.cache.DeletePetDetail(ctx, userID)
	}

	// 3. 发布领域事件（事务外，异步处理）
	if s.publisher != nil && len(events) > 0 {
		for _, event := range events {
			if e, ok := event.(shared.Event); ok {
				_ = s.publisher.Publish(ctx, e)
			}
		}
	}

	return response, nil
}

// CreatePet 创建宠物
func (s *Service) CreatePet(ctx context.Context, userID int64, req CreatePetRequest) (*CreatePetResponse, error) {
	var dto *PetDetailDTO

	err := s.uow.Do(ctx, func(txCtx context.Context) error {

		// 使用领域服务创建宠物
		var p *pet.Pet
		var err error
		// TODO
		req.SpeciesID = "" // TODO 临时处理
		if req.SpeciesID != "" {
			// 指定物种创建，解析物种ID
			speciesIDInt, parseErr := strconv.Atoi(req.SpeciesID)
			if parseErr != nil {
				return ErrInvalidSpeciesID
			}
			p, err = s.petDomainSvc.CreatePet(userID, req.Name, pet.SpeciesID(speciesIDInt))
			if err != nil {
				return err
			}
		} else {
			// 随机物种创建
			p = s.petDomainSvc.CreateRandomPet(userID, req.Name)
		}

		// 保存
		if err := s.petRepo.Save(txCtx, p); err != nil {
			return err
		}

		dto = s.toPetDetailDTO(p)
		return nil
	})

	if err != nil {
		return nil, err
	}

	return &CreatePetResponse{Pet: *dto}, nil
}

// PlayWithPet 和宠物玩耍
func (s *Service) PlayWithPet(ctx context.Context, userID int64) (*PlayPetResponse, error) {
	var response *PlayPetResponse

	err := s.uow.Do(ctx, func(txCtx context.Context) error {
		p, err := s.petRepo.FindByUserID(txCtx, userID)
		if err != nil {
			return err
		}

		oldLevel := p.Level

		if err := p.Play(); err != nil {
			return err
		}

		if err := s.petRepo.Save(txCtx, p); err != nil {
			return err
		}

		response = &PlayPetResponse{
			Happiness: p.Happiness,
			Energy:    p.Energy,
			ExpGained: 10,
			LevelUp:   p.Level > oldLevel,
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// 清除缓存
	if s.cache != nil {
		_ = s.cache.DeletePetDetail(ctx, userID)
	}

	return response, nil
}

// CleanPet 清洁宠物
func (s *Service) CleanPet(ctx context.Context, userID int64) (*CleanPetResponse, error) {
	var response *CleanPetResponse

	err := s.uow.Do(ctx, func(txCtx context.Context) error {
		p, err := s.petRepo.FindByUserID(txCtx, userID)
		if err != nil {
			return err
		}

		if err := p.Clean(); err != nil {
			return err
		}

		if err := s.petRepo.Save(txCtx, p); err != nil {
			return err
		}

		response = &CleanPetResponse{
			Cleanliness: p.Cleanliness,
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	if s.cache != nil {
		_ = s.cache.DeletePetDetail(ctx, userID)
	}

	return response, nil
}

// --- DTO 转换 ---

func (s *Service) toPetDetailDTO(p *pet.Pet) *PetDetailDTO {
	return &PetDetailDTO{
		ID:   p.ID,
		Name: p.Name,
		Appearance: AppearanceDTO{
			ColorPrimary:   p.Appearance.ColorPrimary,
			ColorSecondary: p.Appearance.ColorSecondary,
			PatternType:    p.Appearance.PatternTypeName(),
			BodyType:       p.Appearance.BodyTypeName(),
			Description:    p.Appearance.String(),
		},
		Personality: PersonalityDTO{
			Activity:    p.Personality.Activity,
			Appetite:    p.Personality.Appetite,
			Social:      p.Personality.Social,
			Curiosity:   p.Personality.Curiosity,
			Description: p.Personality.Describe(),
		},
		Skill: SkillDTO{
			Name:        p.Skill.Name(),
			Level:       p.Skill.Level,
			Rarity:      p.Skill.Rarity(),
			Description: p.Skill.Description(),
		},
		Stage:     p.StageName(),
		Level:     p.Level,
		Exp:       p.Exp,
		ExpToNext: p.Level * 100,
		Status: StatusDTO{
			Hunger:      p.Hunger,
			Happiness:   p.Happiness,
			Cleanliness: p.Cleanliness,
			Energy:      p.Energy,
			IsHungry:    p.IsHungry(),
			IsUnhappy:   p.IsUnhappy(),
			IsDirty:     p.IsDirty(),
			IsTired:     p.IsTired(),
		},
		GeneCode: p.Gene.String(),
	}
}

// ============================================================
// 繁殖相关方法
// ============================================================

// BreedPets 繁殖宠物
func (s *Service) BreedPets(ctx context.Context, userID int64, req BreedPetsRequest) (*BreedPetsResponse, error) {
	var response *BreedPetsResponse
	var events []any

	err := s.uow.Do(ctx, func(txCtx context.Context) error {
		// 1. 获取父母1
		parent1, err := s.petRepo.FindByID(txCtx, req.Parent1ID)
		if err != nil {
			return err
		}
		// 验证所有权
		if parent1.UserID != userID {
			return ErrNotPetOwner
		}

		var parent2 *pet.Pet
		var result *pet.BreedingResult

		if req.Parent2ID > 0 {
			// 2. 双亲繁殖
			parent2, err = s.petRepo.FindByID(txCtx, req.Parent2ID)
			if err != nil {
				return err
			}
			// 验证所有权
			if parent2.UserID != userID {
				return ErrNotPetOwner
			}

			// 3. 委托领域服务执行繁殖
			result, err = s.petDomainSvc.BreedPets(parent1, parent2, req.ChildName, userID)
			if err != nil {
				return err
			}
		} else {
			// 分裂繁殖
			result, err = s.petDomainSvc.SelfBreedPet(parent1, req.ChildName, userID)
			if err != nil {
				return err
			}
		}

		// 4. 保存后代
		if err := s.petRepo.Save(txCtx, result.Child); err != nil {
			return err
		}

		// 5. 保存父母（繁殖时间已更新）
		if err := s.petRepo.Save(txCtx, parent1); err != nil {
			return err
		}
		if parent2 != nil {
			if err := s.petRepo.Save(txCtx, parent2); err != nil {
				return err
			}
		}

		// 6. 收集事件
		events = append(events, result.Child.Events()...)

		// 7. 构建响应
		inheritedGenes := []string{}
		mutations := []string{}
		if result.IsHidden {
			mutations = append(mutations, "触发隐藏物种融合")
		}

		response = &BreedPetsResponse{
			Offspring:      *s.toPetDetailDTO(result.Child),
			InheritedGenes: inheritedGenes,
			Mutations:      mutations,
			Parent1Updated: PetBreedingStatusDTO{
				ID: parent1.ID,
			},
		}
		if parent2 != nil {
			response.Parent2Updated = &PetBreedingStatusDTO{
				ID: parent2.ID,
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// 清除缓存
	if s.cache != nil {
		_ = s.cache.DeletePetDetail(ctx, userID)
	}

	// 发布事件
	if s.publisher != nil && len(events) > 0 {
		for _, event := range events {
			if e, ok := event.(shared.Event); ok {
				_ = s.publisher.Publish(ctx, e)
			}
		}
	}

	return response, nil
}

// CanBreed 检查是否可以繁殖
func (s *Service) CanBreed(ctx context.Context, userID int64, parent1ID, parent2ID int64) (*CanBreedResponse, error) {
	parent1, err := s.petRepo.FindByID(ctx, parent1ID)
	if err != nil {
		return nil, err
	}
	if parent1.UserID != userID {
		return &CanBreedResponse{CanBreed: false, Reason: "非宠物主人"}, nil
	}

	if parent2ID > 0 {
		parent2, err := s.petRepo.FindByID(ctx, parent2ID)
		if err != nil {
			return nil, err
		}
		if parent2.UserID != userID {
			return &CanBreedResponse{CanBreed: false, Reason: "非宠物主人"}, nil
		}

		// 委托领域服务检查
		if err := s.petDomainSvc.CanBreedPair(parent1, parent2); err != nil {
			return &CanBreedResponse{CanBreed: false, Reason: err.Error()}, nil
		}
	} else {
		// 检查分裂繁殖
		if err := s.petDomainSvc.CanSelfBreed(parent1); err != nil {
			return &CanBreedResponse{CanBreed: false, Reason: err.Error()}, nil
		}
	}

	return &CanBreedResponse{CanBreed: true}, nil
}

// PredictOffspring 预测后代物种
func (s *Service) PredictOffspring(ctx context.Context, req PredictOffspringRequest) (*PredictOffspringResponse, error) {
	parent1, err := s.petRepo.FindByID(ctx, req.Parent1ID)
	if err != nil {
		return nil, err
	}

	var parent2 *pet.Pet
	if req.Parent2ID > 0 {
		parent2, err = s.petRepo.FindByID(ctx, req.Parent2ID)
		if err != nil {
			return nil, err
		}
	}

	// 委托领域服务预测
	predictions := s.petDomainSvc.PredictOffspringSpecies(parent1, parent2)

	// 转换为 DTO
	result := make([]SpeciesProbabilityDTO, 0, len(predictions))
	for _, p := range predictions {
		species, ok := s.petDomainSvc.GetSpecies(p.SpeciesID)
		name := strconv.Itoa(int(p.SpeciesID))
		if ok {
			name = species.Name
		}
		result = append(result, SpeciesProbabilityDTO{
			SpeciesID:   strconv.Itoa(int(p.SpeciesID)),
			SpeciesName: name,
			Probability: float64(p.Probability) / 100.0, // 转换为 0-1 范围
		})
	}

	return &PredictOffspringResponse{PossibleSpecies: result}, nil
}

// ============================================================
// 评分和物种相关方法
// ============================================================

// GetPetScore 获取宠物评分
func (s *Service) GetPetScore(ctx context.Context, userID int64) (*GetPetScoreResponse, error) {
	p, err := s.petRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// 委托领域服务计算评分
	score := s.petDomainSvc.CalculatePetScore(p)

	// 获取物种稀有度
	rarityScore := 0
	if species, ok := s.petDomainSvc.GetSpecies(p.SpeciesID); ok {
		rarityScore = species.Rarity * 100
	}

	return &GetPetScoreResponse{
		Score: score,
		Breakdown: ScoreBreakdown{
			LevelScore:      p.Level * 10,
			SkillScore:      p.Skill.Strength * 50,
			RarityScore:     rarityScore,
			StatusScore:     (p.Hunger + p.Happiness + p.Cleanliness) / 3,
			StageScore:      int(p.Stage) * 100,
			GenerationScore: p.Generation * 20,
		},
	}, nil
}

// GetAvailableSpecies 获取可用物种列表
func (s *Service) GetAvailableSpecies(ctx context.Context) (*GetSpeciesListResponse, error) {
	species := s.petDomainSvc.GetAvailableSpecies()

	result := make([]SpeciesDTO, 0, len(species))
	for _, sp := range species {
		result = append(result, SpeciesDTO{
			ID:       strconv.Itoa(int(sp.ID)),
			Name:     sp.Name,
			Category: sp.Category.Name(),
			Rarity:   sp.Rarity,
		})
	}

	return &GetSpeciesListResponse{Species: result}, nil
}

// 应用层错误
var (
	ErrPetNotFound      = errors.New("宠物不存在")
	ErrAlreadyHasPet    = errors.New("已经拥有宠物了")
	ErrInvalidFoodItem  = errors.New("无效的食物道具")
	ErrNotPetOwner      = errors.New("非宠物主人")
	ErrInvalidSpeciesID = errors.New("无效的物种ID")
)
