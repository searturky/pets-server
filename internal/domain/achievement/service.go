// Package achievement 成就领域
// 领域服务 - 成就检查逻辑
package achievement

import "context"

// DomainService 成就领域服务
type DomainService struct {
	repo Repository
}

// NewDomainService 创建领域服务
func NewDomainService(repo Repository) *DomainService {
	return &DomainService{repo: repo}
}

// CheckAndUnlock 检查并解锁成就
// 返回新解锁的成就列表
func (s *DomainService) CheckAndUnlock(
	ctx context.Context,
	userID int64,
	conditionType ConditionType,
	currentValue int,
) ([]*UserAchievement, error) {
	// 获取所有成就定义
	definitions, err := s.repo.GetAllDefinitions(ctx)
	if err != nil {
		return nil, err
	}

	var unlocked []*UserAchievement

	for _, def := range definitions {
		// 检查条件类型是否匹配
		if def.ConditionType != conditionType {
			continue
		}

		// 检查是否达成条件
		if currentValue < def.ConditionValue {
			continue
		}

		// 检查是否已获得
		has, err := s.repo.HasAchievement(ctx, userID, def.ID)
		if err != nil {
			return nil, err
		}
		if has {
			continue
		}

		// 创建并保存用户成就
		ua := NewUserAchievement(userID, def.ID)
		if err := s.repo.Save(ctx, ua); err != nil {
			return nil, err
		}

		unlocked = append(unlocked, ua)
	}

	return unlocked, nil
}

