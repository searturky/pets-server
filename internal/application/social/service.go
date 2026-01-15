// Package social 社交应用服务
// 处理社交相关的业务用例
package social

import (
	"context"
	"errors"

	"pets-server/internal/domain/shared"
	"pets-server/internal/domain/social"
)

// Service 社交应用服务
type Service struct {
	friendRepo social.FriendRepository
	giftRepo   social.GiftRepository
	tradeRepo  social.TradeRepository
	visitRepo  social.VisitRepository
	uow        shared.UnitOfWork
	publisher  shared.EventPublisher
}

// NewService 创建社交应用服务
func NewService(
	friendRepo social.FriendRepository,
	giftRepo social.GiftRepository,
	tradeRepo social.TradeRepository,
	visitRepo social.VisitRepository,
	uow shared.UnitOfWork,
	publisher shared.EventPublisher,
) *Service {
	return &Service{
		friendRepo: friendRepo,
		giftRepo:   giftRepo,
		tradeRepo:  tradeRepo,
		visitRepo:  visitRepo,
		uow:        uow,
		publisher:  publisher,
	}
}

// --- 好友功能 ---

// SendFriendRequest 发送好友申请
func (s *Service) SendFriendRequest(ctx context.Context, userID int, req AddFriendRequest) error {
	if userID == req.FriendID {
		return social.ErrCannotAddSelf
	}

	return s.uow.Do(ctx, func(txCtx context.Context) error {
		// 检查是否已经是好友
		existing, err := s.friendRepo.FindByUsers(txCtx, userID, req.FriendID)
		if err != nil && !errors.Is(err, social.ErrFriendshipNotFound) {
			return err
		}
		if existing != nil {
			if existing.IsFriend() {
				return social.ErrAlreadyFriends
			}
			// 如果有待处理的申请，不重复创建
			return nil
		}

		// 创建好友申请
		friendship := social.NewFriendRequest(userID, req.FriendID)
		return s.friendRepo.Save(txCtx, friendship)
	})
}

// AcceptFriendRequest 接受好友申请
func (s *Service) AcceptFriendRequest(ctx context.Context, userID int, friendshipID int) error {
	return s.uow.Do(ctx, func(txCtx context.Context) error {
		friendship, err := s.friendRepo.FindByID(txCtx, friendshipID)
		if err != nil {
			return err
		}

		// 验证是接收者才能接受
		if friendship.FriendID != userID {
			return ErrNotAllowed
		}

		if err := friendship.Accept(); err != nil {
			return err
		}

		return s.friendRepo.Save(txCtx, friendship)
	})
}

// GetFriendList 获取好友列表
func (s *Service) GetFriendList(ctx context.Context, userID int) (*FriendListResponse, error) {
	friendships, err := s.friendRepo.FindFriends(ctx, userID)
	if err != nil {
		return nil, err
	}

	friends := make([]FriendDTO, 0, len(friendships))
	for _, f := range friendships {
		// TODO: 批量获取用户信息和宠物信息
		friends = append(friends, FriendDTO{
			UserID:        f.FriendID,
			Intimacy:      f.Intimacy,
			IntimacyLevel: f.IntimacyLevel(),
		})
	}

	return &FriendListResponse{
		Friends: friends,
		Total:   len(friends),
	}, nil
}

// --- 礼物功能 ---

// SendGift 发送礼物
func (s *Service) SendGift(ctx context.Context, userID int, req SendGiftRequest) error {
	// TODO: 完整实现需要：
	// 1. 验证是好友关系
	// 2. 扣除发送者道具
	// 3. 增加接收者道具
	// 4. 创建礼物记录
	// 5. 增加亲密度
	// 6. 发布事件

	return s.uow.Do(ctx, func(txCtx context.Context) error {
		// 创建礼物记录
		gift := social.NewGiftRecord(userID, req.ToUserID, req.ItemID, req.Quantity, req.Message)
		return s.giftRepo.Save(txCtx, gift)
	})
}

// GetReceivedGifts 获取收到的礼物
func (s *Service) GetReceivedGifts(ctx context.Context, userID int) ([]GiftRecordDTO, error) {
	gifts, err := s.giftRepo.FindByReceiver(ctx, userID, false)
	if err != nil {
		return nil, err
	}

	result := make([]GiftRecordDTO, 0, len(gifts))
	for _, g := range gifts {
		result = append(result, GiftRecordDTO{
			ID:         g.ID,
			FromUserID: g.FromUserID,
			ItemID:     g.ItemID,
			Quantity:   g.Quantity,
			Message:    g.Message,
			CreatedAt:  g.CreatedAt,
			IsRead:     g.IsRead,
		})
	}

	return result, nil
}

// 应用层错误
var (
	ErrNotAllowed = errors.New("无权进行此操作")
)
