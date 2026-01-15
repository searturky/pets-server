// Package social 社交领域
// Repository 仓储接口
package social

import "context"

// FriendRepository 好友仓储接口
type FriendRepository interface {
	// FindByID 根据ID查找好友关系
	FindByID(ctx context.Context, id int) (*Friendship, error)

	// FindByUsers 根据两个用户ID查找好友关系
	FindByUsers(ctx context.Context, userID, friendID int) (*Friendship, error)

	// FindFriends 获取用户的所有好友
	FindFriends(ctx context.Context, userID int) ([]*Friendship, error)

	// FindPendingRequests 获取待处理的好友申请
	FindPendingRequests(ctx context.Context, userID int) ([]*Friendship, error)

	// Save 保存好友关系
	Save(ctx context.Context, friendship *Friendship) error

	// Delete 删除好友关系
	Delete(ctx context.Context, id int) error
}

// GiftRepository 礼物仓储接口
type GiftRepository interface {
	// FindByID 根据ID查找礼物记录
	FindByID(ctx context.Context, id int) (*GiftRecord, error)

	// FindByReceiver 获取用户收到的礼物
	FindByReceiver(ctx context.Context, userID int, onlyUnread bool) ([]*GiftRecord, error)

	// FindBySender 获取用户发送的礼物
	FindBySender(ctx context.Context, userID int) ([]*GiftRecord, error)

	// Save 保存礼物记录
	Save(ctx context.Context, gift *GiftRecord) error
}

// TradeRepository 交易仓储接口
type TradeRepository interface {
	// FindByID 根据ID查找交易
	FindByID(ctx context.Context, id int) (*Trade, error)

	// FindByUser 获取用户相关的交易
	FindByUser(ctx context.Context, userID int) ([]*Trade, error)

	// FindPending 获取待处理的交易
	FindPending(ctx context.Context, userID int) ([]*Trade, error)

	// Save 保存交易
	Save(ctx context.Context, trade *Trade) error
}

// VisitRepository 拜访记录仓储接口
type VisitRepository interface {
	// RecordVisit 记录拜访
	RecordVisit(ctx context.Context, visitorID, hostID int) error

	// CountTodayVisits 统计今日被拜访次数
	CountTodayVisits(ctx context.Context, hostID int) (int, error)

	// HasVisitedToday 今天是否已拜访过
	HasVisitedToday(ctx context.Context, visitorID, hostID int) (bool, error)
}

