// Package social 社交领域
// Gift 礼物实体
package social

import "time"

// GiftRecord 礼物记录实体
type GiftRecord struct {
	ID         int64
	FromUserID int64  // 送礼者
	ToUserID   int64  // 收礼者
	ItemID     int    // 道具ID
	Quantity   int    // 数量
	Message    string // 留言
	IsRead     bool   // 是否已读
	CreatedAt  time.Time
}

// NewGiftRecord 创建礼物记录
func NewGiftRecord(fromUserID, toUserID int64, itemID, quantity int, message string) *GiftRecord {
	return &GiftRecord{
		FromUserID: fromUserID,
		ToUserID:   toUserID,
		ItemID:     itemID,
		Quantity:   quantity,
		Message:    message,
		IsRead:     false,
		CreatedAt:  time.Now(),
	}
}

// MarkAsRead 标记为已读
func (g *GiftRecord) MarkAsRead() {
	g.IsRead = true
}

// GiftSentEvent 礼物发送事件
type GiftSentEvent struct {
	GiftID     int64     `json:"gift_id"`
	FromUserID int64     `json:"from_user_id"`
	ToUserID   int64     `json:"to_user_id"`
	ItemID     int       `json:"item_id"`
	Quantity   int       `json:"quantity"`
	Timestamp  time.Time `json:"timestamp"`
}

func (e GiftSentEvent) EventName() string { return "social.gift_sent" }

