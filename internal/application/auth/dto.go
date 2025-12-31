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
	ID        int64  `json:"id"`
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

