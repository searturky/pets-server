// Package auth 认证应用服务
// 处理用户认证相关的用例
package auth

import (
	"context"
	"errors"
	"time"

	"pets-server/internal/domain/shared"
	"pets-server/internal/domain/user"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// Service 认证应用服务
type Service struct {
	userRepo       user.Repository
	uow            shared.UnitOfWork
	wechatAuthFunc func(code string) (openID string, err error) // 微信认证函数
	jwtSecret      string
	jwtExpireHours int
	sessionStore   SessionStore
}

// NewService 创建认证服务
func NewService(
	userRepo user.Repository,
	uow shared.UnitOfWork,
	wechatAuthFunc func(code string) (string, error),
	jwtSecret string,
	jwtExpireHours int,
	sessionStore SessionStore,
) *Service {
	return &Service{
		userRepo:       userRepo,
		uow:            uow,
		wechatAuthFunc: wechatAuthFunc,
		jwtSecret:      jwtSecret,
		jwtExpireHours: jwtExpireHours,
		sessionStore:   sessionStore,
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
	token, err := s.issueTokenWithSession(ctx, u.ID)
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

// Login 账号密码登录
func (s *Service) Login(ctx context.Context, req LoginRequest) (*LoginResponse, error) {
	// 1. 根据用户名查找用户
	u, err := s.userRepo.FindByUsername(ctx, req.Username)
	if err != nil {
		if errors.Is(err, user.ErrUserNotFound) {
			return nil, errors.New("用户名或密码错误")
		}
		return nil, err
	}

	// 2. 验证密码
	if err := u.VerifyPassword(req.Password); err != nil {
		return nil, errors.New("用户名或密码错误")
	}

	// 3. 更新登录时间
	err = s.uow.Do(ctx, func(txCtx context.Context) error {
		u.UpdateLogin()
		return s.userRepo.Save(txCtx, u)
	})
	if err != nil {
		return nil, err
	}

	// 4. 生成 JWT Token
	token, err := s.issueTokenWithSession(ctx, u.ID)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		Token: token,
		UserInfo: UserInfo{
			ID:        u.ID,
			Nickname:  u.Nickname,
			AvatarURL: u.AvatarURL,
			Coins:     u.Coins,
			Diamonds:  u.Diamonds,
		},
	}, nil
}

// Register 账号注册
func (s *Service) Register(ctx context.Context, req RegisterRequest) (*LoginResponse, error) {
	// 1. 检查用户名是否已存在
	existing, err := s.userRepo.FindByUsername(ctx, req.Username)
	if err != nil && !errors.Is(err, user.ErrUserNotFound) {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("用户名已存在")
	}

	// 2. 创建新用户
	u, err := user.NewUserWithPassword(req.Username, req.Password, req.Nickname)
	if err != nil {
		return nil, errors.New("创建用户失败: " + err.Error())
	}

	// 3. 保存用户
	err = s.uow.Do(ctx, func(txCtx context.Context) error {
		return s.userRepo.Save(txCtx, u)
	})
	if err != nil {
		return nil, err
	}

	// 4. 生成 JWT Token
	token, err := s.issueTokenWithSession(ctx, u.ID)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		Token: token,
		UserInfo: UserInfo{
			ID:        u.ID,
			Nickname:  u.Nickname,
			AvatarURL: u.AvatarURL,
			Coins:     u.Coins,
			Diamonds:  u.Diamonds,
		},
	}, nil
}

// GetUserInfo 获取用户信息
func (s *Service) GetUserInfo(ctx context.Context, userID int) (*UserInfo, error) {
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

// Logout 退出登录（清除当前会话）
func (s *Service) Logout(ctx context.Context, userID int) error {
	if s.sessionStore == nil {
		return errors.New("session store not configured")
	}
	return s.sessionStore.DeleteCurrentSession(ctx, userID)
}

// KickUser 主动踢下线（清除指定用户会话）
func (s *Service) KickUser(ctx context.Context, userID int) error {
	if s.sessionStore == nil {
		return errors.New("session store not configured")
	}
	return s.sessionStore.DeleteCurrentSession(ctx, userID)
}

// issueTokenWithSession 生成会话并返回 JWT Token
func (s *Service) issueTokenWithSession(ctx context.Context, userID int) (string, error) {
	if s.sessionStore == nil {
		return "", errors.New("session store not configured")
	}

	sessionID := uuid.NewString()
	token, err := s.generateToken(userID, sessionID)
	if err != nil {
		return "", err
	}

	ttl := time.Duration(s.jwtExpireHours) * time.Hour
	if err := s.sessionStore.SetCurrentSession(ctx, userID, sessionID, ttl); err != nil {
		return "", err
	}

	return token, nil
}

// generateToken 生成 JWT Token
func (s *Service) generateToken(userID int, sessionID string) (string, error) {
	now := time.Now()
	expiresAt := now.Add(time.Duration(s.jwtExpireHours) * time.Hour)

	// 创建 Claims
	claims := jwt.MapClaims{
		"user_id":    userID,
		"session_id": sessionID,
		"iat":        now.Unix(),       // 签发时间
		"exp":        expiresAt.Unix(), // 过期时间
		"nbf":        now.Unix(),       // 生效时间
	}

	// 创建 Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 签名并获取完整的 token 字符串
	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", errors.New("生成 token 失败: " + err.Error())
	}

	return tokenString, nil
}
