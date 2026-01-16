// Package model GORM 模型定义
package model

import "time"

// Friendship 好友关系表
type Friendship struct {
	BaseModel
	UserID      int       `gorm:"column:user_id;index;not null;comment:用户ID"`
	FriendID    int       `gorm:"column:friend_id;index;not null;comment:好友ID"`
	Status      int16     `gorm:"default:0;comment:状态(0待确认1已通过2已拒绝)"` // 0待确认 1已通过 2已拒绝
	Intimacy    int       `gorm:"default:0;comment:亲密度(0-100)"`        // 亲密度 0-100
	ConfirmedAt time.Time `gorm:"column:confirmed_at;comment:确认时间"`
}

// TableName 表名
func (Friendship) TableName() string {
	return "friendships"
}

// GiftRecord 礼物记录表
type GiftRecord struct {
	BaseModel
	FromUserID int    `gorm:"column:from_user_id;index;not null;comment:送礼用户ID"`
	ToUserID   int    `gorm:"column:to_user_id;index;not null;comment:收礼用户ID"`
	ItemID     int    `gorm:"column:item_id;not null;comment:道具ID"`
	Quantity   int    `gorm:"default:1;comment:数量"`
	Message    string `gorm:"type:varchar(128);comment:留言"`
	IsRead     bool   `gorm:"column:is_read;default:false;comment:是否已读"`
}

// TableName 表名
func (GiftRecord) TableName() string {
	return "gift_records"
}

// Trade 交易表
type Trade struct {
	BaseModel
	FromUserID      int       `gorm:"column:from_user_id;index;not null;comment:发起方用户ID"`
	ToUserID        int       `gorm:"column:to_user_id;index;not null;comment:接收方用户ID"`
	OfferItemID     int       `gorm:"column:offer_item_id;comment:提供的道具ID"`
	OfferQuantity   int       `gorm:"column:offer_quantity;comment:提供的数量"`
	RequestItemID   int       `gorm:"column:request_item_id;comment:请求的道具ID"`
	RequestQuantity int       `gorm:"column:request_quantity;comment:请求的数量"`
	Status          int16     `gorm:"default:0;comment:状态(0待确认1已完成2已取消3已过期)"` // 0待确认 1已完成 2已取消 3已过期
	CompletedAt     time.Time `gorm:"column:completed_at;comment:完成时间"`
}

// TableName 表名
func (Trade) TableName() string {
	return "trades"
}

// VisitRecord 访问记录表
type VisitRecord struct {
	BaseModel
	VisitorID int `gorm:"column:visitor_id;index;not null;comment:访客用户ID"`
	HostID    int `gorm:"column:host_id;index;not null;comment:主人用户ID"`
}

// TableName 表名
func (VisitRecord) TableName() string {
	return "visit_records"
}
