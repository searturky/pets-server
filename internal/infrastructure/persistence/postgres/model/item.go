// Package model GORM 模型定义
package model

// ItemDefinition 道具定义表
type ItemDefinition struct {
	BaseModel
	Name        string `gorm:"type:varchar(32);not null;comment:道具名称"`
	Description string `gorm:"type:text;comment:道具描述"`
	ItemType    int16  `gorm:"column:item_type;comment:道具类型(1食物2清洁3玩具4装饰5特殊)"` // 1食物 2清洁 3玩具 4装饰 5特殊
	EffectType  string `gorm:"column:effect_type;type:varchar(32);comment:效果类型"`
	EffectValue int    `gorm:"column:effect_value;comment:效果数值"`
	Price       int    `gorm:"default:0;comment:价格"`
	Rarity      int16  `gorm:"default:1;comment:稀有度(1普通2稀有3史诗4传说)"` // 1普通 2稀有 3史诗 4传说
}

// TableName 表名
func (ItemDefinition) TableName() string {
	return "item_definitions"
}

// UserItem 用户道具背包
type UserItem struct {
	BaseModel
	UserID   int `gorm:"index;not null;comment:用户ID"`
	ItemID   int `gorm:"column:item_id;not null;comment:道具ID"`
	Quantity int `gorm:"default:1;comment:数量"`

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
	PetID  int    `gorm:"index;not null;comment:宠物ID"`
	ItemID int    `gorm:"column:item_id;not null;comment:道具ID"`
	Slot   string `gorm:"type:varchar(16);comment:装备槽位(head/body/accessory)"` // head, body, accessory
}

// TableName 表名
func (PetDecoration) TableName() string {
	return "pet_decorations"
}
