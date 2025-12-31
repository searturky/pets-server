// Package social 社交领域
// 领域服务
package social

// TODO: 社交领域服务
// - 好友推荐
// - 交易撮合
// - 社交排行

// DomainService 社交领域服务
type DomainService struct {
	friendRepo FriendRepository
	giftRepo   GiftRepository
	tradeRepo  TradeRepository
}

// NewDomainService 创建领域服务
func NewDomainService(
	friendRepo FriendRepository,
	giftRepo GiftRepository,
	tradeRepo TradeRepository,
) *DomainService {
	return &DomainService{
		friendRepo: friendRepo,
		giftRepo:   giftRepo,
		tradeRepo:  tradeRepo,
	}
}

