// Package shared 包含跨领域的共享接口和类型
// UnitOfWork (工作单元) 接口定义
// 用于管理跨多个仓储的事务一致性
package shared

import "context"

// UnitOfWork 工作单元接口
// 定义在领域层，由基础设施层实现
// 用于保证多个仓储操作的事务一致性
type UnitOfWork interface {
	// Do 在事务中执行操作
	// 如果 fn 返回 error，事务将回滚
	// 如果 fn 返回 nil，事务将提交
	Do(ctx context.Context, fn func(ctx context.Context) error) error
}

