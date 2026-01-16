// Package model GORM 模型定义
package model

import "time"

// User 用户表
type User struct {
	BaseModel
	Username    string    `gorm:"type:varchar(32);uniqueIndex;comment:用户名"`
	Password    string    `gorm:"type:varchar(128);comment:密码哈希"`
	OpenID      string    `gorm:"column:open_id;type:varchar(64);uniqueIndex;comment:微信OpenID"`
	UnionID     string    `gorm:"column:union_id;type:varchar(64);comment:微信UnionID"`
	Nickname    string    `gorm:"type:varchar(32);not null;comment:昵称"`
	AvatarURL   string    `gorm:"column:avatar_url;type:varchar(256);comment:头像URL"`
	Coins       int       `gorm:"default:0;comment:金币数量"`
	Diamonds    int       `gorm:"default:0;comment:钻石数量"`
	LastLoginAt time.Time `gorm:"column:last_login_at;comment:最后登录时间"`
}

// TableName 表名
func (User) TableName() string {
	return "users"
}
