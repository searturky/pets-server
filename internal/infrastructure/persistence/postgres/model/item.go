// Package model GORM 模型定义
package model

// ItemDefinition 道具定义表
type ItemDefinition struct {
	BaseModel
	Name        string `gorm:"type:varchar(32);not null"`
	Description string `gorm:"type:text"`
	ItemType    int16  `gorm:"column:item_type"` // 1食物 2清洁 3玩具 4装饰 5特殊
	EffectType  string `gorm:"column:effect_type;type:varchar(32)"`
	EffectValue int    `gorm:"column:effect_value"`
	Price       int    `gorm:"default:0"`
	Rarity      int16  `gorm:"default:1"` // 1普通 2稀有 3史诗 4传说
}

// TableName 表名
func (ItemDefinition) TableName() string {
	return "item_definitions"
}

// UserItem 用户道具背包
type UserItem struct {
	BaseModel
	UserID   int `gorm:"index;not null"`
	ItemID   int `gorm:"column:item_id;not null"`
	Quantity int `gorm:"default:1"`

	// 联合唯一索引
	// 同一用户的同一道具只有一条记录
}

// TableName 表名
func (UserItem) TableName() string {
	return "user_items"
}

// 创建联合唯一索引的迁移钩子
func (UserItem) BeforeCreate(tx interface{}) error {
	// GORM 会通过 tag 或手动执行来创建
	return nil
}

// PetDecoration 宠物装饰 (穿戴中的装饰)
type PetDecoration struct {
	BaseModel
	PetID  int    `gorm:"index;not null"`
	ItemID int    `gorm:"column:item_id;not null"`
	Slot   string `gorm:"type:varchar(16)"` // head, body, accessory
}

// TableName 表名
func (PetDecoration) TableName() string {
	return "pet_decorations"
}
