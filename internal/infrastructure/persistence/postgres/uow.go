// Package postgres UnitOfWork 实现
package postgres

import (
	"context"

	"gorm.io/gorm"
)

// 事务上下文键
type txKey struct{}

// UnitOfWork GORM 实现的工作单元
type UnitOfWork struct {
	db *gorm.DB
}

// NewUnitOfWork 创建工作单元
func NewUnitOfWork(db *gorm.DB) *UnitOfWork {
	return &UnitOfWork{db: db}
}

// Do 在事务中执行操作
// 如果 fn 返回 error，事务回滚
// 如果 fn 返回 nil，事务提交
func (u *UnitOfWork) Do(ctx context.Context, fn func(ctx context.Context) error) error {
	return u.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 将事务 tx 放入 context
		txCtx := context.WithValue(ctx, txKey{}, tx)
		return fn(txCtx)
	})
}

// GetTx 从 context 获取当前事务
// 如果没有事务，返回默认 db
func GetTx(ctx context.Context, defaultDB *gorm.DB) *gorm.DB {
	if tx, ok := ctx.Value(txKey{}).(*gorm.DB); ok {
		return tx
	}
	return defaultDB.WithContext(ctx)
}

