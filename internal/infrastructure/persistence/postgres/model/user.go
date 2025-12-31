// Package model GORM 模型定义
package model

import "time"

// User 用户表
type User struct {
	ID          int64     `gorm:"primaryKey;autoIncrement"`
	OpenID      string    `gorm:"column:open_id;type:varchar(64);uniqueIndex"`
	UnionID     string    `gorm:"column:union_id;type:varchar(64)"`
	Nickname    string    `gorm:"type:varchar(32);not null"`
	AvatarURL   string    `gorm:"column:avatar_url;type:varchar(256)"`
	Coins       int       `gorm:"default:0"`
	Diamonds    int       `gorm:"default:0"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	LastLoginAt time.Time `gorm:"column:last_login_at"`
}

// TableName 表名
func (User) TableName() string {
	return "users"
}

