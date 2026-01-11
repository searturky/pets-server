// Package repo 仓储实现
package repo

import (
	"context"
	"encoding/json"
	"errors"

	"gorm.io/gorm"

	"pets-server/internal/domain/pet"
	"pets-server/internal/infrastructure/persistence/postgres"
	"pets-server/internal/infrastructure/persistence/postgres/model"
)

// PetRepository 宠物仓储实现
type PetRepository struct {
	db *gorm.DB
}

// NewPetRepository 创建宠物仓储
func NewPetRepository(db *gorm.DB) *PetRepository {
	return &PetRepository{db: db}
}

// FindByID 根据ID查找宠物
func (r *PetRepository) FindByID(ctx context.Context, id int64) (*pet.Pet, error) {
	db := postgres.GetTx(ctx, r.db)

	var m model.Pet
	if err := db.First(&m, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, pet.ErrPetNotFound
		}
		return nil, err
	}

	return r.toDomain(&m), nil
}

// FindByUserID 根据用户ID查找宠物
func (r *PetRepository) FindByUserID(ctx context.Context, userID int64) (*pet.Pet, error) {
	db := postgres.GetTx(ctx, r.db)

	var m model.Pet
	if err := db.Where("user_id = ?", userID).First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, pet.ErrPetNotFound
		}
		return nil, err
	}

	return r.toDomain(&m), nil
}

// FindByUserIDAll 根据用户ID查找所有宠物
func (r *PetRepository) FindByUserIDAll(ctx context.Context, userID int64) ([]*pet.Pet, error) {
	db := postgres.GetTx(ctx, r.db)

	var models []model.Pet
	if err := db.Where("user_id = ?", userID).Find(&models).Error; err != nil {
		return nil, err
	}

	pets := make([]*pet.Pet, len(models))
	for i, m := range models {
		pets[i] = r.toDomain(&m)
	}

	return pets, nil
}

// Save 保存宠物
func (r *PetRepository) Save(ctx context.Context, p *pet.Pet) error {
	db := postgres.GetTx(ctx, r.db)

	m := r.toModel(p)
	if err := db.Save(m).Error; err != nil {
		return err
	}

	// 回写ID
	p.ID = m.ID
	return nil
}

// Delete 删除宠物
func (r *PetRepository) Delete(ctx context.Context, id int64) error {
	db := postgres.GetTx(ctx, r.db)
	return db.Delete(&model.Pet{}, id).Error
}

// FindAll 查找所有宠物（分页）
func (r *PetRepository) FindAll(ctx context.Context, offset, limit int) ([]*pet.Pet, error) {
	db := postgres.GetTx(ctx, r.db)

	var models []model.Pet
	if err := db.Offset(offset).Limit(limit).Find(&models).Error; err != nil {
		return nil, err
	}

	pets := make([]*pet.Pet, len(models))
	for i, m := range models {
		pets[i] = r.toDomain(&m)
	}

	return pets, nil
}

// FindBySpecies 根据物种查找宠物
func (r *PetRepository) FindBySpecies(ctx context.Context, speciesID int, offset, limit int) ([]*pet.Pet, error) {
	db := postgres.GetTx(ctx, r.db)

	var models []model.Pet
	if err := db.Where("species_id = ?", speciesID).Offset(offset).Limit(limit).Find(&models).Error; err != nil {
		return nil, err
	}

	pets := make([]*pet.Pet, len(models))
	for i, m := range models {
		pets[i] = r.toDomain(&m)
	}

	return pets, nil
}

// FindByParent 根据父母ID查找子代
func (r *PetRepository) FindByParent(ctx context.Context, parentID int64) ([]*pet.Pet, error) {
	db := postgres.GetTx(ctx, r.db)

	var models []model.Pet
	if err := db.Where("parent1_id = ? OR parent2_id = ?", parentID, parentID).Find(&models).Error; err != nil {
		return nil, err
	}

	pets := make([]*pet.Pet, len(models))
	for i, m := range models {
		pets[i] = r.toDomain(&m)
	}

	return pets, nil
}

// CountAll 统计宠物总数
func (r *PetRepository) CountAll(ctx context.Context) (int64, error) {
	db := postgres.GetTx(ctx, r.db)

	var count int64
	if err := db.Model(&model.Pet{}).Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

// --- 模型转换 ---

func (r *PetRepository) toDomain(m *model.Pet) *pet.Pet {
	gene := pet.NewGene(m.GeneCode)

	// 解析物种特有外观
	var specialAppearance pet.SpecialAppearance
	if m.SpecialAppearance != "" {
		json.Unmarshal([]byte(m.SpecialAppearance), &specialAppearance)
	} else {
		specialAppearance = pet.NewSpecialAppearance()
	}

	p := &pet.Pet{
		ID:        m.ID,
		UserID:    m.UserID,
		Name:      m.Name,
		SpeciesID: pet.SpeciesID(m.SpeciesID),
		Gender:    pet.Gender(m.Gender),
		Gene:      gene,
		Appearance: pet.Appearance{
			ColorPrimary:   m.ColorPrimary,
			ColorSecondary: m.ColorSecondary,
			PatternType:    int(m.PatternType),
			PatternDensity: int(m.PatternDensity),
			BodyType:       int(m.BodyType),
			EyeShape:       int(m.EyeShape),
			EyeColor:       int(m.EyeColor),
		},
		SpecialAppearance: specialAppearance,
		Personality: pet.Personality{
			Activity:     int(m.TraitActivity),
			Appetite:     int(m.TraitAppetite),
			Social:       int(m.TraitSocial),
			Curiosity:    int(m.TraitCuriosity),
			Temper:       int(m.TraitTemper),
			Loyalty:      int(m.TraitLoyalty),
			Intelligence: int(m.TraitIntelligence),
			Playfulness:  int(m.TraitPlayfulness),
		},
		Skill: pet.Skill{
			Type:     pet.SkillType(m.SkillID),
			Level:    int(m.SkillLevel),
			Strength: int(m.SkillStrength),
		},
		Stage:         pet.Stage(m.Stage),
		Exp:           m.Exp,
		Level:         m.Level,
		Hunger:        int(m.Hunger),
		Happiness:     int(m.Happiness),
		Cleanliness:   int(m.Cleanliness),
		Energy:        int(m.Energy),
		Parent1ID:     m.Parent1ID,
		Parent2ID:     m.Parent2ID,
		Generation:    m.Generation,
		LastBreedAt:   m.LastBreedAt,
		LastFedAt:     m.LastFedAt,
		LastPlayedAt:  m.LastPlayedAt,
		LastCleanedAt: m.LastCleanedAt,
		BornAt:        m.BornAt,
		CreatedAt:     m.CreatedAt,
	}

	return p
}

func (r *PetRepository) toModel(p *pet.Pet) *model.Pet {
	// 序列化物种特有外观
	specialAppearanceJSON, _ := json.Marshal(p.SpecialAppearance)

	return &model.Pet{
		ID:                p.ID,
		UserID:            p.UserID,
		Name:              p.Name,
		SpeciesID:         int(p.SpeciesID),
		Gender:            int16(p.Gender),
		GeneCode:          p.Gene.String(),
		ColorPrimary:      p.Appearance.ColorPrimary,
		ColorSecondary:    p.Appearance.ColorSecondary,
		PatternType:       int16(p.Appearance.PatternType),
		PatternDensity:    int16(p.Appearance.PatternDensity),
		BodyType:          int16(p.Appearance.BodyType),
		EyeShape:          int16(p.Appearance.EyeShape),
		EyeColor:          int16(p.Appearance.EyeColor),
		SpecialAppearance: string(specialAppearanceJSON),
		TraitActivity:     int16(p.Personality.Activity),
		TraitAppetite:     int16(p.Personality.Appetite),
		TraitSocial:       int16(p.Personality.Social),
		TraitCuriosity:    int16(p.Personality.Curiosity),
		TraitTemper:       int16(p.Personality.Temper),
		TraitLoyalty:      int16(p.Personality.Loyalty),
		TraitIntelligence: int16(p.Personality.Intelligence),
		TraitPlayfulness:  int16(p.Personality.Playfulness),
		SkillID:           int(p.Skill.Type),
		SkillLevel:        int16(p.Skill.Level),
		SkillStrength:     int16(p.Skill.Strength),
		SkillSecondaryID:  0, // TODO: 支持副技能
		Stage:             int16(p.Stage),
		Exp:               p.Exp,
		Level:             p.Level,
		Hunger:            int16(p.Hunger),
		Happiness:         int16(p.Happiness),
		Cleanliness:       int16(p.Cleanliness),
		Energy:            int16(p.Energy),
		Parent1ID:         p.Parent1ID,
		Parent2ID:         p.Parent2ID,
		Generation:        p.Generation,
		LastBreedAt:       p.LastBreedAt,
		LastFedAt:         p.LastFedAt,
		LastPlayedAt:      p.LastPlayedAt,
		LastCleanedAt:     p.LastCleanedAt,
		BornAt:            p.BornAt,
		CreatedAt:         p.CreatedAt,
	}
}
