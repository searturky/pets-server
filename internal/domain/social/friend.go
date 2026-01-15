// Package social 社交领域
// Friend 好友关系实体
package social

import (
	"errors"
	"time"
)

// FriendStatus 好友状态
type FriendStatus int

const (
	FriendStatusPending  FriendStatus = 0 // 待确认
	FriendStatusAccepted FriendStatus = 1 // 已通过
	FriendStatusRejected FriendStatus = 2 // 已拒绝
)

// Friendship 好友关系实体
type Friendship struct {
	ID          int
	UserID      int          // 发起者
	FriendID    int          // 接收者
	Status      FriendStatus // 状态
	Intimacy    int          // 亲密度 0-100
	CreatedAt   time.Time    // 申请时间
	ConfirmedAt time.Time    // 确认时间
}

// NewFriendRequest 创建好友申请
func NewFriendRequest(userID, friendID int) *Friendship {
	return &Friendship{
		UserID:    userID,
		FriendID:  friendID,
		Status:    FriendStatusPending,
		Intimacy:  0,
		CreatedAt: time.Now(),
	}
}

// Accept 接受好友申请
func (f *Friendship) Accept() error {
	if f.Status != FriendStatusPending {
		return ErrInvalidFriendStatus
	}
	f.Status = FriendStatusAccepted
	f.Intimacy = 10 // 初始亲密度
	f.ConfirmedAt = time.Now()
	return nil
}

// Reject 拒绝好友申请
func (f *Friendship) Reject() error {
	if f.Status != FriendStatusPending {
		return ErrInvalidFriendStatus
	}
	f.Status = FriendStatusRejected
	f.ConfirmedAt = time.Now()
	return nil
}

// AddIntimacy 增加亲密度
func (f *Friendship) AddIntimacy(amount int) {
	if f.Status != FriendStatusAccepted {
		return
	}
	f.Intimacy = min(f.Intimacy+amount, 100)
}

// IntimacyLevel 亲密度等级
func (f *Friendship) IntimacyLevel() string {
	switch {
	case f.Intimacy >= 80:
		return "挚友"
	case f.Intimacy >= 60:
		return "好友"
	case f.Intimacy >= 40:
		return "朋友"
	case f.Intimacy >= 20:
		return "熟人"
	default:
		return "点头之交"
	}
}

// IsFriend 是否是好友
func (f *Friendship) IsFriend() bool {
	return f.Status == FriendStatusAccepted
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// 领域错误
var (
	ErrInvalidFriendStatus = errors.New("无效的好友状态")
	ErrAlreadyFriends      = errors.New("已经是好友了")
	ErrFriendshipNotFound  = errors.New("好友关系不存在")
	ErrCannotAddSelf       = errors.New("不能添加自己为好友")
)

