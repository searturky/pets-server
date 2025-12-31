// Package user 用户领域
// 领域服务 - 处理不属于单个实体的业务逻辑
package user

// TODO: 领域服务
// 当业务逻辑不属于单个实体时，放在领域服务中
// 例如：
// - 用户合并逻辑
// - 跨用户的业务规则验证

// DomainService 用户领域服务
type DomainService struct {
	repo Repository
}

// NewDomainService 创建领域服务
func NewDomainService(repo Repository) *DomainService {
	return &DomainService{repo: repo}
}

