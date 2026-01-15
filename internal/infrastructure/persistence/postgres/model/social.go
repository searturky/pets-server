// Package model GORM 模型定义
package model

import "time"

// Friendship 好友关系表
type Friendship struct {
	BaseModel
	UserID      int `gorm:"column:user_id;index;not null"`
	FriendID    int `gorm:"column:friend_id;index;not null"`
	Status      int16     `gorm:"default:0"` // 0待确认 1已通过 2已拒绝
	Intimacy    int       `gorm:"default:0"` // 亲密度 0-100
	ConfirmedAt time.Time `gorm:"column:confirmed_at"`
}

// TableName 表名
func (Friendship) TableName() string {
	return "friendships"
}

// GiftRecord 礼物记录表
type GiftRecord struct {
	BaseModel
	FromUserID int `gorm:"column:from_user_id;index;not null"`
	ToUserID   int `gorm:"column:to_user_id;index;not null"`
	ItemID     int    `gorm:"column:item_id;not null"`
	Quantity   int    `gorm:"default:1"`
	Message    string `gorm:"type:varchar(128)"`
	IsRead     bool   `gorm:"column:is_read;default:false"`
}

// TableName 表名
func (GiftRecord) TableName() string {
	return "gift_records"
}

// Trade 交易表
type Trade struct {
	BaseModel
	FromUserID      int `gorm:"column:from_user_id;index;not null"`
	ToUserID        int `gorm:"column:to_user_id;index;not null"`
	OfferItemID     int       `gorm:"column:offer_item_id"`
	OfferQuantity   int       `gorm:"column:offer_quantity"`
	RequestItemID   int       `gorm:"column:request_item_id"`
	RequestQuantity int       `gorm:"column:request_quantity"`
	Status          int16     `gorm:"default:0"` // 0待确认 1已完成 2已取消 3已过期
	CompletedAt     time.Time `gorm:"column:completed_at"`
}

// TableName 表名
func (Trade) TableName() string {
	return "trades"
}

// VisitRecord 访问记录表
type VisitRecord struct {
	BaseModel
	VisitorID int `gorm:"column:visitor_id;index;not null"`
	HostID    int `gorm:"column:host_id;index;not null"`
}

// TableName 表名
func (VisitRecord) TableName() string {
	return "visit_records"
}

