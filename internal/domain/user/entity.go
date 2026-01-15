// Package user 用户领域
// 包含用户实体、值对象和领域逻辑
package user

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

// User 用户实体
// 聚合根，代表游戏中的玩家
type User struct {
	ID          int       // 用户唯一标识
	Username    string    // 用户名（用于账号密码登录）
	Password    string    // 密码哈希（用于账号密码登录）
	OpenID      string    // 微信 OpenID
	UnionID     string    // 微信 UnionID
	Nickname    string    // 昵称
	AvatarURL   string    // 头像URL
	Coins       int       // 金币（普通货币）
	Diamonds    int       // 钻石（高级货币）
	CreatedAt   time.Time // 创建时间
	LastLoginAt time.Time // 最后登录时间
}

// NewUser 创建新用户
func NewUser(openID, nickname, avatarURL string) *User {
	now := time.Now()
	return &User{
		OpenID:      openID,
		Nickname:    nickname,
		AvatarURL:   avatarURL,
		Coins:       100, // 初始金币
		Diamonds:    10,  // 初始钻石
		CreatedAt:   now,
		LastLoginAt: now,
	}
}

// AddCoins 增加金币
func (u *User) AddCoins(amount int) {
	if amount > 0 {
		u.Coins += amount
	}
}

// SpendCoins 消费金币
func (u *User) SpendCoins(amount int) error {
	if amount <= 0 {
		return ErrInvalidAmount
	}
	if u.Coins < amount {
		return ErrInsufficientCoins
	}
	u.Coins -= amount
	return nil
}

// AddDiamonds 增加钻石
func (u *User) AddDiamonds(amount int) {
	if amount > 0 {
		u.Diamonds += amount
	}
}

// SpendDiamonds 消费钻石
func (u *User) SpendDiamonds(amount int) error {
	if amount <= 0 {
		return ErrInvalidAmount
	}
	if u.Diamonds < amount {
		return ErrInsufficientDiamonds
	}
	u.Diamonds -= amount
	return nil
}

// UpdateLogin 更新登录时间
func (u *User) UpdateLogin() {
	u.LastLoginAt = time.Now()
}

// NewUserWithPassword 创建新用户（使用账号密码）
func NewUserWithPassword(username, password, nickname string) (*User, error) {
	// 密码加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	return &User{
		Username:    username,
		Password:    string(hashedPassword),
		Nickname:    nickname,
		Coins:       100, // 初始金币
		Diamonds:    10,  // 初始钻石
		CreatedAt:   now,
		LastLoginAt: now,
	}, nil
}

// VerifyPassword 验证密码
func (u *User) VerifyPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}

