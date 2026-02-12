// Package repo 仓储实现
package repo

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"pets-server/internal/domain/user"
	"pets-server/internal/infrastructure/persistence/postgres"
	"pets-server/internal/infrastructure/persistence/postgres/model"
)

// UserRepository 用户仓储实现
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository 创建用户仓储
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// FindByID 根据ID查找用户
func (r *UserRepository) FindByID(ctx context.Context, id int) (*user.User, error) {
	db := postgres.GetTx(ctx, r.db)

	var m model.User
	if err := db.First(&m, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, user.ErrUserNotFound
		}
		return nil, err
	}

	return r.toDomain(&m), nil
}

// FindByUsername 根据用户名查找用户
func (r *UserRepository) FindByUsername(ctx context.Context, username string) (*user.User, error) {
	db := postgres.GetTx(ctx, r.db)

	var m model.User
	if err := db.Where("username = ?", username).First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, user.ErrUserNotFound
		}
		return nil, err
	}

	return r.toDomain(&m), nil
}

// FindByOpenID 根据微信OpenID查找用户
func (r *UserRepository) FindByOpenID(ctx context.Context, openID string) (*user.User, error) {
	db := postgres.GetTx(ctx, r.db)

	var m model.User
	if err := db.Where("open_id = ?", openID).First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, user.ErrUserNotFound
		}
		return nil, err
	}

	return r.toDomain(&m), nil
}

// Save 保存用户
func (r *UserRepository) Save(ctx context.Context, u *user.User) error {
	db := postgres.GetTx(ctx, r.db)

	m := r.toModel(u)
	if err := db.Save(m).Error; err != nil {
		return err
	}

	// 回写ID
	u.ID = m.ID
	return nil
}

// Delete 删除用户
func (r *UserRepository) Delete(ctx context.Context, id int) error {
	db := postgres.GetTx(ctx, r.db)
	return db.Delete(&model.User{}, id).Error
}

// --- 模型转换 ---

func (r *UserRepository) toDomain(m *model.User) *user.User {
	return &user.User{
		ID:          m.ID,
		Username:    m.Username,
		Password:    m.Password,
		OpenID:      m.OpenID,
		UnionID:     m.UnionID,
		Nickname:    m.Nickname,
		AvatarURL:   m.AvatarURL,
		Coins:       m.Coins,
		Diamonds:    m.Diamonds,
		ActivePetID: m.ActivePetID,
		CreatedAt:   m.CreatedAt,
		LastLoginAt: m.LastLoginAt,
	}
}

func (r *UserRepository) toModel(u *user.User) *model.User {
	m := &model.User{
		Username:    u.Username,
		Password:    u.Password,
		OpenID:      u.OpenID,
		UnionID:     u.UnionID,
		Nickname:    u.Nickname,
		AvatarURL:   u.AvatarURL,
		Coins:       u.Coins,
		Diamonds:    u.Diamonds,
		ActivePetID: u.ActivePetID,
		LastLoginAt: u.LastLoginAt,
	}
	m.ID = u.ID
	return m
}
