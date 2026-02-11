// Package auth 认证应用服务
// DTO 数据传输对象
package auth

// WxLoginRequest 微信登录请求
type WxLoginRequest struct {
	Code string `json:"code" binding:"required"` // 微信登录凭证
}

// WxLoginResponse 微信登录响应
type WxLoginResponse struct {
	Token    string   `json:"token"`    // JWT Token
	UserInfo UserInfo `json:"userInfo"` // 用户信息
	IsNew    bool     `json:"isNew"`    // 是否新用户
}

// UserInfo 用户信息
type UserInfo struct {
	ID        int    `json:"id"`
	Nickname  string `json:"nickname"`
	AvatarURL string `json:"avatarUrl"`
	Coins     int    `json:"coins"`
	Diamonds  int    `json:"diamonds"`
}

// UpdateProfileRequest 更新用户信息请求
type UpdateProfileRequest struct {
	Nickname  string `json:"nickname"`
	AvatarURL string `json:"avatarUrl"`
}

// LoginRequest 账号密码登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required,min=6,max=20,alphanum"` // 用户名
	Password string `json:"password" binding:"required,min=6,max=32"`          // 密码
}

// RegisterRequest 账号注册请求
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=6,max=20,alphanum"` // 用户名，只能包含字母、数字和下划线
	Password string `json:"password" binding:"required,min=6,max=32"`          // 密码
	Nickname string `json:"nickname" binding:"required,min=2,max=20"`          // 昵称
}

// KickRequest 踢下线请求
type KickRequest struct {
	UserID int `json:"userId" binding:"required"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	Token    string   `json:"token"`    // JWT Token
	UserInfo UserInfo `json:"userInfo"` // 用户信息
}
