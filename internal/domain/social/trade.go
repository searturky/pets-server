// Package social 社交领域
// Trade 交易实体
package social

import (
	"errors"
	"time"
)

// TradeStatus 交易状态
type TradeStatus int

const (
	TradeStatusPending   TradeStatus = 0 // 待确认
	TradeStatusCompleted TradeStatus = 1 // 已完成
	TradeStatusCancelled TradeStatus = 2 // 已取消
	TradeStatusExpired   TradeStatus = 3 // 已过期
)

// Trade 交易实体
type Trade struct {
	ID              int64
	FromUserID      int64       // 发起者
	ToUserID        int64       // 接收者
	OfferItemID     int         // 提供的道具ID
	OfferQuantity   int         // 提供的数量
	RequestItemID   int         // 请求的道具ID
	RequestQuantity int         // 请求的数量
	Status          TradeStatus // 状态
	CreatedAt       time.Time
	CompletedAt     time.Time
}

// NewTrade 创建交易
func NewTrade(fromUserID, toUserID int64, offerItemID, offerQty, requestItemID, requestQty int) *Trade {
	return &Trade{
		FromUserID:      fromUserID,
		ToUserID:        toUserID,
		OfferItemID:     offerItemID,
		OfferQuantity:   offerQty,
		RequestItemID:   requestItemID,
		RequestQuantity: requestQty,
		Status:          TradeStatusPending,
		CreatedAt:       time.Now(),
	}
}

// Accept 接受交易
func (t *Trade) Accept() error {
	if t.Status != TradeStatusPending {
		return ErrInvalidTradeStatus
	}
	t.Status = TradeStatusCompleted
	t.CompletedAt = time.Now()
	return nil
}

// Cancel 取消交易
func (t *Trade) Cancel() error {
	if t.Status != TradeStatusPending {
		return ErrInvalidTradeStatus
	}
	t.Status = TradeStatusCancelled
	t.CompletedAt = time.Now()
	return nil
}

// IsExpired 检查是否过期（24小时）
func (t *Trade) IsExpired() bool {
	if t.Status != TradeStatusPending {
		return false
	}
	return time.Since(t.CreatedAt) > 24*time.Hour
}

// MarkExpired 标记为过期
func (t *Trade) MarkExpired() {
	if t.Status == TradeStatusPending {
		t.Status = TradeStatusExpired
		t.CompletedAt = time.Now()
	}
}

// 领域错误
var (
	ErrInvalidTradeStatus = errors.New("无效的交易状态")
	ErrTradeNotFound      = errors.New("交易不存在")
	ErrTradeExpired       = errors.New("交易已过期")
)

