// Package social 社交应用服务
// DTO 数据传输对象
package social

import "time"

// --- 好友相关 ---

// AddFriendRequest 添加好友请求
type AddFriendRequest struct {
	FriendID int64 `json:"friendId" binding:"required"`
}

// FriendDTO 好友DTO
type FriendDTO struct {
	UserID        int64  `json:"userId"`
	Nickname      string `json:"nickname"`
	AvatarURL     string `json:"avatarUrl"`
	Intimacy      int    `json:"intimacy"`
	IntimacyLevel string `json:"intimacyLevel"`
	PetName       string `json:"petName,omitempty"`
	PetLevel      int    `json:"petLevel,omitempty"`
}

// FriendRequestDTO 好友申请DTO
type FriendRequestDTO struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"userId"`
	Nickname  string    `json:"nickname"`
	AvatarURL string    `json:"avatarUrl"`
	CreatedAt time.Time `json:"createdAt"`
}

// FriendListResponse 好友列表响应
type FriendListResponse struct {
	Friends []FriendDTO `json:"friends"`
	Total   int         `json:"total"`
}

// --- 礼物相关 ---

// SendGiftRequest 发送礼物请求
type SendGiftRequest struct {
	ToUserID int64  `json:"toUserId" binding:"required"`
	ItemID   int    `json:"itemId" binding:"required"`
	Quantity int    `json:"quantity" binding:"required,min=1"`
	Message  string `json:"message" binding:"max=128"`
}

// GiftRecordDTO 礼物记录DTO
type GiftRecordDTO struct {
	ID           int64     `json:"id"`
	FromUserID   int64     `json:"fromUserId"`
	FromNickname string    `json:"fromNickname"`
	ItemID       int       `json:"itemId"`
	ItemName     string    `json:"itemName"`
	Quantity     int       `json:"quantity"`
	Message      string    `json:"message"`
	CreatedAt    time.Time `json:"createdAt"`
	IsRead       bool      `json:"isRead"`
}

// --- 交易相关 ---

// CreateTradeRequest 创建交易请求
type CreateTradeRequest struct {
	ToUserID        int64 `json:"toUserId" binding:"required"`
	OfferItemID     int   `json:"offerItemId" binding:"required"`
	OfferQuantity   int   `json:"offerQuantity" binding:"required,min=1"`
	RequestItemID   int   `json:"requestItemId" binding:"required"`
	RequestQuantity int   `json:"requestQuantity" binding:"required,min=1"`
}

// TradeDTO 交易DTO
type TradeDTO struct {
	ID              int64     `json:"id"`
	FromUserID      int64     `json:"fromUserId"`
	FromNickname    string    `json:"fromNickname"`
	ToUserID        int64     `json:"toUserId"`
	ToNickname      string    `json:"toNickname"`
	OfferItemName   string    `json:"offerItemName"`
	OfferQuantity   int       `json:"offerQuantity"`
	RequestItemName string    `json:"requestItemName"`
	RequestQuantity int       `json:"requestQuantity"`
	Status          string    `json:"status"`
	CreatedAt       time.Time `json:"createdAt"`
}

// --- 拜访相关 ---

// VisitResponse 拜访响应
type VisitResponse struct {
	Pet          PetVisitDTO `json:"pet"`
	RewardCoins  int         `json:"rewardCoins"`
	CanVisitMore bool        `json:"canVisitMore"`
}

// PetVisitDTO 拜访时的宠物信息
type PetVisitDTO struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Level     int    `json:"level"`
	Stage     string `json:"stage"`
	OwnerName string `json:"ownerName"`
}

