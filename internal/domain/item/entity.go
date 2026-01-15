// Package item 道具领域
// 包含道具实体、背包和领域逻辑
package item

import (
	"errors"
	"time"
)

// ItemType 道具类型
type ItemType int

const (
	ItemTypeFood       ItemType = 1 // 食物
	ItemTypeClean      ItemType = 2 // 清洁用品
	ItemTypeToy        ItemType = 3 // 玩具
	ItemTypeDecoration ItemType = 4 // 装饰品
	ItemTypeSpecial    ItemType = 5 // 特殊道具
)

// ItemDefinition 道具定义（值对象）
// 定义道具的基本属性，不可变
type ItemDefinition struct {
	ID          int
	Name        string
	Description string
	Type        ItemType
	EffectType  string // 效果类型
	EffectValue int    // 效果数值
	Price       int    // 购买价格
	Rarity      int    // 稀有度 1-4
}

// UserItem 用户道具（实体）
// 表示用户背包中的道具
type UserItem struct {
	ID        int
	UserID    int
	ItemID    int      // 对应 ItemDefinition.ID
	Quantity  int      // 数量
	CreatedAt time.Time
}

// NewUserItem 创建用户道具
func NewUserItem(userID int, itemID int, quantity int) *UserItem {
	return &UserItem{
		UserID:    userID,
		ItemID:    itemID,
		Quantity:  quantity,
		CreatedAt: time.Now(),
	}
}

// Add 增加数量
func (i *UserItem) Add(amount int) {
	if amount > 0 {
		i.Quantity += amount
	}
}

// Consume 消耗道具
func (i *UserItem) Consume(amount int) error {
	if amount <= 0 {
		return ErrInvalidAmount
	}
	if i.Quantity < amount {
		return ErrInsufficientItem
	}
	i.Quantity -= amount
	return nil
}

// IsEmpty 是否已用完
func (i *UserItem) IsEmpty() bool {
	return i.Quantity <= 0
}

// PetDecoration 宠物装饰（实体）
// 宠物当前穿戴的装饰品
type PetDecoration struct {
	ID         int
	PetID      int
	ItemID     int    // 装饰品道具ID
	Slot       string // 槽位: head, body, accessory
	EquippedAt time.Time
}

// 领域错误
var (
	ErrInvalidAmount    = errors.New("无效的数量")
	ErrInsufficientItem = errors.New("道具数量不足")
	ErrItemNotFound     = errors.New("道具不存在")
)

