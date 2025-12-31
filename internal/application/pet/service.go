// Package pet 宠物应用服务
// 编排宠物相关的业务用例
// 这是完整的读写示例
package pet

import (
	"context"
	"errors"
	"time"

	"pets-server/internal/domain/item"
	"pets-server/internal/domain/pet"
	"pets-server/internal/domain/shared"
)

// Service 宠物应用服务
type Service struct {
	petRepo   pet.Repository
	itemRepo  item.Repository
	uow       shared.UnitOfWork
	publisher shared.EventPublisher
	cache     CacheService // 缓存服务接口
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
	uow shared.UnitOfWork,
	publisher shared.EventPublisher,
	cache CacheService,
) *Service {
	return &Service{
		petRepo:   petRepo,
		itemRepo:  itemRepo,
		uow:       uow,
		publisher: publisher,
		cache:     cache,
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
		// 检查用户是否已有宠物
		existing, err := s.petRepo.FindByUserID(txCtx, userID)
		if err != nil && !errors.Is(err, pet.ErrPetNotFound) {
			return err
		}
		if existing != nil {
			return ErrAlreadyHasPet
		}

		// 创建新宠物（随机基因）
		p := pet.NewPet(userID, req.Name)

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

// 应用层错误
var (
	ErrPetNotFound     = errors.New("宠物不存在")
	ErrAlreadyHasPet   = errors.New("已经拥有宠物了")
	ErrInvalidFoodItem = errors.New("无效的食物道具")
)

