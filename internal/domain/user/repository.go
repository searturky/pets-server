// Package user 用户领域
// Repository 接口定义 - 用户仓储
package user

import "context"

// Repository 用户仓储接口
// 定义在领域层，由基础设施层实现
type Repository interface {
	// FindByID 根据ID查找用户
	FindByID(ctx context.Context, id int64) (*User, error)

	// FindByOpenID 根据微信OpenID查找用户
	FindByOpenID(ctx context.Context, openID string) (*User, error)

	// Save 保存用户（新增或更新）
	Save(ctx context.Context, user *User) error

	// Delete 删除用户
	Delete(ctx context.Context, id int64) error
}

