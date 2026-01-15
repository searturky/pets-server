// Package repo 仓储实现
package repo

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"pets-server/internal/domain/item"
	"pets-server/internal/infrastructure/persistence/postgres"
	"pets-server/internal/infrastructure/persistence/postgres/model"
)

// ItemRepository 道具仓储实现
type ItemRepository struct {
	db *gorm.DB
}

// NewItemRepository 创建道具仓储
func NewItemRepository(db *gorm.DB) *ItemRepository {
	return &ItemRepository{db: db}
}

// FindByID 根据ID查找用户道具
func (r *ItemRepository) FindByID(ctx context.Context, id int) (*item.UserItem, error) {
	db := postgres.GetTx(ctx, r.db)

	var m model.UserItem
	if err := db.First(&m, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, item.ErrItemNotFound
		}
		return nil, err
	}

	return r.toDomain(&m), nil
}

// FindByUserAndItem 根据用户ID和道具ID查找
func (r *ItemRepository) FindByUserAndItem(ctx context.Context, userID int, itemID int) (*item.UserItem, error) {
	db := postgres.GetTx(ctx, r.db)

	var m model.UserItem
	if err := db.Where("user_id = ? AND item_id = ?", userID, itemID).First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, item.ErrItemNotFound
		}
		return nil, err
	}

	return r.toDomain(&m), nil
}

// FindByUserID 获取用户所有道具
func (r *ItemRepository) FindByUserID(ctx context.Context, userID int) ([]*item.UserItem, error) {
	db := postgres.GetTx(ctx, r.db)

	var models []model.UserItem
	if err := db.Where("user_id = ?", userID).Find(&models).Error; err != nil {
		return nil, err
	}

	items := make([]*item.UserItem, len(models))
	for i, m := range models {
		items[i] = r.toDomain(&m)
	}

	return items, nil
}

// FindByUserAndType 获取用户某类型的所有道具
func (r *ItemRepository) FindByUserAndType(ctx context.Context, userID int, itemType item.ItemType) ([]*item.UserItem, error) {
	db := postgres.GetTx(ctx, r.db)

	var models []model.UserItem
	if err := db.Joins("JOIN item_definitions ON user_items.item_id = item_definitions.id").
		Where("user_items.user_id = ? AND item_definitions.item_type = ?", userID, itemType).
		Find(&models).Error; err != nil {
		return nil, err
	}

	items := make([]*item.UserItem, len(models))
	for i, m := range models {
		items[i] = r.toDomain(&m)
	}

	return items, nil
}

// Save 保存用户道具
func (r *ItemRepository) Save(ctx context.Context, i *item.UserItem) error {
	db := postgres.GetTx(ctx, r.db)

	m := r.toModel(i)
	if err := db.Save(m).Error; err != nil {
		return err
	}

	i.ID = m.ID
	return nil
}

// Delete 删除用户道具
func (r *ItemRepository) Delete(ctx context.Context, id int) error {
	db := postgres.GetTx(ctx, r.db)
	return db.Delete(&model.UserItem{}, id).Error
}

// GetDefinition 获取道具定义
func (r *ItemRepository) GetDefinition(ctx context.Context, itemID int) (*item.ItemDefinition, error) {
	db := postgres.GetTx(ctx, r.db)

	var m model.ItemDefinition
	if err := db.First(&m, itemID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, item.ErrItemNotFound
		}
		return nil, err
	}

	return &item.ItemDefinition{
		ID:          int(m.ID),
		Name:        m.Name,
		Description: m.Description,
		Type:        item.ItemType(m.ItemType),
		EffectType:  m.EffectType,
		EffectValue: m.EffectValue,
		Price:       m.Price,
		Rarity:      int(m.Rarity),
	}, nil
}

// GetAllDefinitions 获取所有道具定义
func (r *ItemRepository) GetAllDefinitions(ctx context.Context) ([]*item.ItemDefinition, error) {
	db := postgres.GetTx(ctx, r.db)

	var models []model.ItemDefinition
	if err := db.Find(&models).Error; err != nil {
		return nil, err
	}

	defs := make([]*item.ItemDefinition, len(models))
	for i, m := range models {
		defs[i] = &item.ItemDefinition{
			ID:          int(m.ID),
			Name:        m.Name,
			Description: m.Description,
			Type:        item.ItemType(m.ItemType),
			EffectType:  m.EffectType,
			EffectValue: m.EffectValue,
			Price:       m.Price,
			Rarity:      int(m.Rarity),
		}
	}

	return defs, nil
}

// --- 模型转换 ---

func (r *ItemRepository) toDomain(m *model.UserItem) *item.UserItem {
	return &item.UserItem{
		ID:        m.ID,
		UserID:    m.UserID,
		ItemID:    m.ItemID,
		Quantity:  m.Quantity,
		CreatedAt: m.CreatedAt,
	}
}

func (r *ItemRepository) toModel(i *item.UserItem) *model.UserItem {
	m := &model.UserItem{
		UserID:   i.UserID,
		ItemID:   i.ItemID,
		Quantity: i.Quantity,
	}
	m.ID = i.ID
	return m
}

