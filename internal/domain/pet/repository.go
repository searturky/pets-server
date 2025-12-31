// Package pet 宠物领域
// Repository 仓储接口
package pet

import "context"

// Repository 宠物仓储接口
// 定义在领域层，由基础设施层实现
type Repository interface {
	// FindByID 根据ID查找宠物
	FindByID(ctx context.Context, id int64) (*Pet, error)

	// FindByUserID 根据用户ID查找宠物
	FindByUserID(ctx context.Context, userID int64) (*Pet, error)

	// Save 保存宠物（新增或更新）
	Save(ctx context.Context, pet *Pet) error

	// Delete 删除宠物
	Delete(ctx context.Context, id int64) error

	// FindAll 查找所有宠物（用于定时任务）
	FindAll(ctx context.Context, offset, limit int) ([]*Pet, error)

	// CountAll 统计宠物总数
	CountAll(ctx context.Context) (int64, error)
}

