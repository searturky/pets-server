// Package item 道具领域
// 领域服务
package item

// TODO: 道具领域服务
// - 道具合成
// - 道具交换
// - 背包整理

// DomainService 道具领域服务
type DomainService struct {
	repo Repository
}

// NewDomainService 创建领域服务
func NewDomainService(repo Repository) *DomainService {
	return &DomainService{repo: repo}
}

