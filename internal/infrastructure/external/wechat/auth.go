// Package wechat 微信接口
package wechat

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Config 微信配置
type Config struct {
	AppID     string
	AppSecret string
}

// AuthService 微信认证服务
type AuthService struct {
	appID     string
	appSecret string
}

// NewAuthService 创建微信认证服务
func NewAuthService(cfg Config) *AuthService {
	return &AuthService{
		appID:     cfg.AppID,
		appSecret: cfg.AppSecret,
	}
}

// Code2Session 通过临时登录凭证获取用户信息
// https://developers.weixin.qq.com/miniprogram/dev/api-backend/open-api/login/auth.code2Session.html
func (s *AuthService) Code2Session(code string) (string, error) {
	url := fmt.Sprintf(
		"https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code",
		s.appID, s.appSecret, code,
	)

	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	var result struct {
		OpenID     string `json:"openid"`
		SessionKey string `json:"session_key"`
		UnionID    string `json:"unionid"`
		ErrCode    int    `json:"errcode"`
		ErrMsg     string `json:"errmsg"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("decode response failed: %w", err)
	}

	if result.ErrCode != 0 {
		return "", fmt.Errorf("wechat error: %d - %s", result.ErrCode, result.ErrMsg)
	}

	return result.OpenID, nil
}

// GetOpenID 获取 openid（Code2Session 的简化版）
func (s *AuthService) GetOpenID(code string) (string, error) {
	return s.Code2Session(code)
}

