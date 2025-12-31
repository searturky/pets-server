// Package auth 认证应用服务
// 处理用户认证相关的用例
package auth

import (
	"context"
	"errors"
	"time"

	"pets-server/internal/domain/shared"
	"pets-server/internal/domain/user"
)

// Service 认证应用服务
type Service struct {
	userRepo       user.Repository
	uow            shared.UnitOfWork
	wechatAuthFunc func(code string) (openID string, err error) // 微信认证函数
	jwtSecret      string
	jwtExpireHours int
}

// NewService 创建认证服务
func NewService(
	userRepo user.Repository,
	uow shared.UnitOfWork,
	wechatAuthFunc func(code string) (string, error),
	jwtSecret string,
	jwtExpireHours int,
) *Service {
	return &Service{
		userRepo:       userRepo,
		uow:            uow,
		wechatAuthFunc: wechatAuthFunc,
		jwtSecret:      jwtSecret,
		jwtExpireHours: jwtExpireHours,
	}
}

// WxLogin 微信登录
func (s *Service) WxLogin(ctx context.Context, req WxLoginRequest) (*WxLoginResponse, error) {
	// 1. 调用微信接口获取 openid
	openID, err := s.wechatAuthFunc(req.Code)
	if err != nil {
		return nil, errors.New("微信登录失败: " + err.Error())
	}

	var u *user.User
	var isNew bool

	// 2. 在事务中处理用户
	err = s.uow.Do(ctx, func(txCtx context.Context) error {
		// 查找用户
		existing, err := s.userRepo.FindByOpenID(txCtx, openID)
		if err != nil && !errors.Is(err, user.ErrUserNotFound) {
			return err
		}

		if existing != nil {
			// 已有用户，更新登录时间
			existing.UpdateLogin()
			u = existing
			isNew = false
		} else {
			// 新用户，创建账号
			u = user.NewUser(openID, "新玩家", "")
			isNew = true
		}

		return s.userRepo.Save(txCtx, u)
	})
	if err != nil {
		return nil, err
	}

	// 3. 生成 JWT Token
	token, err := s.generateToken(u.ID)
	if err != nil {
		return nil, err
	}

	return &WxLoginResponse{
		Token: token,
		UserInfo: UserInfo{
			ID:        u.ID,
			Nickname:  u.Nickname,
			AvatarURL: u.AvatarURL,
			Coins:     u.Coins,
			Diamonds:  u.Diamonds,
		},
		IsNew: isNew,
	}, nil
}

// GetUserInfo 获取用户信息
func (s *Service) GetUserInfo(ctx context.Context, userID int64) (*UserInfo, error) {
	u, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &UserInfo{
		ID:        u.ID,
		Nickname:  u.Nickname,
		AvatarURL: u.AvatarURL,
		Coins:     u.Coins,
		Diamonds:  u.Diamonds,
	}, nil
}

// generateToken 生成 JWT Token
func (s *Service) generateToken(userID int64) (string, error) {
	// TODO: 实现 JWT 生成逻辑
	// 这里只是占位，实际应该使用 jwt 库生成
	_ = time.Now().Add(time.Duration(s.jwtExpireHours) * time.Hour)
	return "token_placeholder_" + string(rune(userID)), nil
}

