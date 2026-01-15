// Package model GORM 模型定义
package model

import "time"

// User 用户表
type User struct {
	BaseModel
	Username    string    `gorm:"type:varchar(32);uniqueIndex"`
	Password    string    `gorm:"type:varchar(128)"`
	OpenID      string    `gorm:"column:open_id;type:varchar(64);uniqueIndex"`
	UnionID     string    `gorm:"column:union_id;type:varchar(64)"`
	Nickname    string    `gorm:"type:varchar(32);not null"`
	AvatarURL   string    `gorm:"column:avatar_url;type:varchar(256)"`
	Coins       int       `gorm:"default:0"`
	Diamonds    int       `gorm:"default:0"`
	LastLoginAt time.Time `gorm:"column:last_login_at"`
}

// TableName 表名
func (User) TableName() string {
	return "users"
}
