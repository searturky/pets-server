// Package item 道具领域
// Repository 仓储接口
package item

import "context"

// Repository 道具仓储接口
type Repository interface {
	// --- 用户道具 ---
	
	// FindByID 根据ID查找用户道具
	FindByID(ctx context.Context, id int) (*UserItem, error)

	// FindByUserAndItem 根据用户ID和道具ID查找
	FindByUserAndItem(ctx context.Context, userID int, itemID int) (*UserItem, error)

	// FindByUserID 获取用户所有道具
	FindByUserID(ctx context.Context, userID int) ([]*UserItem, error)

	// FindByUserAndType 获取用户某类型的所有道具
	FindByUserAndType(ctx context.Context, userID int, itemType ItemType) ([]*UserItem, error)

	// Save 保存用户道具
	Save(ctx context.Context, item *UserItem) error

	// Delete 删除用户道具
	Delete(ctx context.Context, id int) error

	// --- 道具定义 ---
	
	// GetDefinition 获取道具定义
	GetDefinition(ctx context.Context, itemID int) (*ItemDefinition, error)

	// GetAllDefinitions 获取所有道具定义
	GetAllDefinitions(ctx context.Context) ([]*ItemDefinition, error)
}

// DecorationRepository 装饰仓储接口
type DecorationRepository interface {
	// FindByPetID 获取宠物的所有装饰
	FindByPetID(ctx context.Context, petID int) ([]*PetDecoration, error)

	// FindByPetAndSlot 获取宠物某槽位的装饰
	FindByPetAndSlot(ctx context.Context, petID int, slot string) (*PetDecoration, error)

	// Save 保存装饰
	Save(ctx context.Context, decoration *PetDecoration) error

	// Delete 删除装饰
	Delete(ctx context.Context, id int) error
}

