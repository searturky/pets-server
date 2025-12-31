// Package repo 仓储实现
package repo

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"

	"pets-server/internal/domain/social"
	"pets-server/internal/infrastructure/persistence/postgres"
	"pets-server/internal/infrastructure/persistence/postgres/model"
)

// FriendRepository 好友仓储实现
type FriendRepository struct {
	db *gorm.DB
}

// NewFriendRepository 创建好友仓储
func NewFriendRepository(db *gorm.DB) *FriendRepository {
	return &FriendRepository{db: db}
}

// FindByID 根据ID查找好友关系
func (r *FriendRepository) FindByID(ctx context.Context, id int64) (*social.Friendship, error) {
	db := postgres.GetTx(ctx, r.db)

	var m model.Friendship
	if err := db.First(&m, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, social.ErrFriendshipNotFound
		}
		return nil, err
	}

	return r.toDomain(&m), nil
}

// FindByUsers 根据两个用户ID查找好友关系
func (r *FriendRepository) FindByUsers(ctx context.Context, userID, friendID int64) (*social.Friendship, error) {
	db := postgres.GetTx(ctx, r.db)

	var m model.Friendship
	if err := db.Where("(user_id = ? AND friend_id = ?) OR (user_id = ? AND friend_id = ?)",
		userID, friendID, friendID, userID).First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, social.ErrFriendshipNotFound
		}
		return nil, err
	}

	return r.toDomain(&m), nil
}

// FindFriends 获取用户的所有好友
func (r *FriendRepository) FindFriends(ctx context.Context, userID int64) ([]*social.Friendship, error) {
	db := postgres.GetTx(ctx, r.db)

	var models []model.Friendship
	if err := db.Where("(user_id = ? OR friend_id = ?) AND status = ?",
		userID, userID, social.FriendStatusAccepted).Find(&models).Error; err != nil {
		return nil, err
	}

	friendships := make([]*social.Friendship, len(models))
	for i, m := range models {
		friendships[i] = r.toDomain(&m)
	}

	return friendships, nil
}

// FindPendingRequests 获取待处理的好友申请
func (r *FriendRepository) FindPendingRequests(ctx context.Context, userID int64) ([]*social.Friendship, error) {
	db := postgres.GetTx(ctx, r.db)

	var models []model.Friendship
	if err := db.Where("friend_id = ? AND status = ?", userID, social.FriendStatusPending).
		Find(&models).Error; err != nil {
		return nil, err
	}

	friendships := make([]*social.Friendship, len(models))
	for i, m := range models {
		friendships[i] = r.toDomain(&m)
	}

	return friendships, nil
}

// Save 保存好友关系
func (r *FriendRepository) Save(ctx context.Context, f *social.Friendship) error {
	db := postgres.GetTx(ctx, r.db)

	m := r.toModel(f)
	if err := db.Save(m).Error; err != nil {
		return err
	}

	f.ID = m.ID
	return nil
}

// Delete 删除好友关系
func (r *FriendRepository) Delete(ctx context.Context, id int64) error {
	db := postgres.GetTx(ctx, r.db)
	return db.Delete(&model.Friendship{}, id).Error
}

func (r *FriendRepository) toDomain(m *model.Friendship) *social.Friendship {
	return &social.Friendship{
		ID:          m.ID,
		UserID:      m.UserID,
		FriendID:    m.FriendID,
		Status:      social.FriendStatus(m.Status),
		Intimacy:    m.Intimacy,
		CreatedAt:   m.CreatedAt,
		ConfirmedAt: m.ConfirmedAt,
	}
}

func (r *FriendRepository) toModel(f *social.Friendship) *model.Friendship {
	return &model.Friendship{
		ID:          f.ID,
		UserID:      f.UserID,
		FriendID:    f.FriendID,
		Status:      int16(f.Status),
		Intimacy:    f.Intimacy,
		CreatedAt:   f.CreatedAt,
		ConfirmedAt: f.ConfirmedAt,
	}
}

// --- GiftRepository ---

// GiftRepository 礼物仓储实现
type GiftRepository struct {
	db *gorm.DB
}

// NewGiftRepository 创建礼物仓储
func NewGiftRepository(db *gorm.DB) *GiftRepository {
	return &GiftRepository{db: db}
}

// FindByID 根据ID查找礼物记录
func (r *GiftRepository) FindByID(ctx context.Context, id int64) (*social.GiftRecord, error) {
	db := postgres.GetTx(ctx, r.db)

	var m model.GiftRecord
	if err := db.First(&m, id).Error; err != nil {
		return nil, err
	}

	return r.toDomain(&m), nil
}

// FindByReceiver 获取用户收到的礼物
func (r *GiftRepository) FindByReceiver(ctx context.Context, userID int64, onlyUnread bool) ([]*social.GiftRecord, error) {
	db := postgres.GetTx(ctx, r.db)

	query := db.Where("to_user_id = ?", userID)
	if onlyUnread {
		query = query.Where("is_read = ?", false)
	}

	var models []model.GiftRecord
	if err := query.Order("created_at DESC").Find(&models).Error; err != nil {
		return nil, err
	}

	gifts := make([]*social.GiftRecord, len(models))
	for i, m := range models {
		gifts[i] = r.toDomain(&m)
	}

	return gifts, nil
}

// FindBySender 获取用户发送的礼物
func (r *GiftRepository) FindBySender(ctx context.Context, userID int64) ([]*social.GiftRecord, error) {
	db := postgres.GetTx(ctx, r.db)

	var models []model.GiftRecord
	if err := db.Where("from_user_id = ?", userID).Order("created_at DESC").Find(&models).Error; err != nil {
		return nil, err
	}

	gifts := make([]*social.GiftRecord, len(models))
	for i, m := range models {
		gifts[i] = r.toDomain(&m)
	}

	return gifts, nil
}

// Save 保存礼物记录
func (r *GiftRepository) Save(ctx context.Context, g *social.GiftRecord) error {
	db := postgres.GetTx(ctx, r.db)

	m := &model.GiftRecord{
		ID:         g.ID,
		FromUserID: g.FromUserID,
		ToUserID:   g.ToUserID,
		ItemID:     g.ItemID,
		Quantity:   g.Quantity,
		Message:    g.Message,
		IsRead:     g.IsRead,
		CreatedAt:  g.CreatedAt,
	}

	if err := db.Save(m).Error; err != nil {
		return err
	}

	g.ID = m.ID
	return nil
}

func (r *GiftRepository) toDomain(m *model.GiftRecord) *social.GiftRecord {
	return &social.GiftRecord{
		ID:         m.ID,
		FromUserID: m.FromUserID,
		ToUserID:   m.ToUserID,
		ItemID:     m.ItemID,
		Quantity:   m.Quantity,
		Message:    m.Message,
		IsRead:     m.IsRead,
		CreatedAt:  m.CreatedAt,
	}
}

// --- TradeRepository ---

// TradeRepository 交易仓储实现
type TradeRepository struct {
	db *gorm.DB
}

// NewTradeRepository 创建交易仓储
func NewTradeRepository(db *gorm.DB) *TradeRepository {
	return &TradeRepository{db: db}
}

// FindByID 根据ID查找交易
func (r *TradeRepository) FindByID(ctx context.Context, id int64) (*social.Trade, error) {
	db := postgres.GetTx(ctx, r.db)

	var m model.Trade
	if err := db.First(&m, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, social.ErrTradeNotFound
		}
		return nil, err
	}

	return r.toDomain(&m), nil
}

// FindByUser 获取用户相关的交易
func (r *TradeRepository) FindByUser(ctx context.Context, userID int64) ([]*social.Trade, error) {
	db := postgres.GetTx(ctx, r.db)

	var models []model.Trade
	if err := db.Where("from_user_id = ? OR to_user_id = ?", userID, userID).
		Order("created_at DESC").Find(&models).Error; err != nil {
		return nil, err
	}

	trades := make([]*social.Trade, len(models))
	for i, m := range models {
		trades[i] = r.toDomain(&m)
	}

	return trades, nil
}

// FindPending 获取待处理的交易
func (r *TradeRepository) FindPending(ctx context.Context, userID int64) ([]*social.Trade, error) {
	db := postgres.GetTx(ctx, r.db)

	var models []model.Trade
	if err := db.Where("to_user_id = ? AND status = ?", userID, social.TradeStatusPending).
		Find(&models).Error; err != nil {
		return nil, err
	}

	trades := make([]*social.Trade, len(models))
	for i, m := range models {
		trades[i] = r.toDomain(&m)
	}

	return trades, nil
}

// Save 保存交易
func (r *TradeRepository) Save(ctx context.Context, t *social.Trade) error {
	db := postgres.GetTx(ctx, r.db)

	m := &model.Trade{
		ID:              t.ID,
		FromUserID:      t.FromUserID,
		ToUserID:        t.ToUserID,
		OfferItemID:     t.OfferItemID,
		OfferQuantity:   t.OfferQuantity,
		RequestItemID:   t.RequestItemID,
		RequestQuantity: t.RequestQuantity,
		Status:          int16(t.Status),
		CreatedAt:       t.CreatedAt,
		CompletedAt:     t.CompletedAt,
	}

	if err := db.Save(m).Error; err != nil {
		return err
	}

	t.ID = m.ID
	return nil
}

func (r *TradeRepository) toDomain(m *model.Trade) *social.Trade {
	return &social.Trade{
		ID:              m.ID,
		FromUserID:      m.FromUserID,
		ToUserID:        m.ToUserID,
		OfferItemID:     m.OfferItemID,
		OfferQuantity:   m.OfferQuantity,
		RequestItemID:   m.RequestItemID,
		RequestQuantity: m.RequestQuantity,
		Status:          social.TradeStatus(m.Status),
		CreatedAt:       m.CreatedAt,
		CompletedAt:     m.CompletedAt,
	}
}

// --- VisitRepository ---

// VisitRepository 拜访记录仓储实现
type VisitRepository struct {
	db *gorm.DB
}

// NewVisitRepository 创建拜访仓储
func NewVisitRepository(db *gorm.DB) *VisitRepository {
	return &VisitRepository{db: db}
}

// RecordVisit 记录拜访
func (r *VisitRepository) RecordVisit(ctx context.Context, visitorID, hostID int64) error {
	db := postgres.GetTx(ctx, r.db)

	m := &model.VisitRecord{
		VisitorID: visitorID,
		HostID:    hostID,
	}

	return db.Create(m).Error
}

// CountTodayVisits 统计今日被拜访次数
func (r *VisitRepository) CountTodayVisits(ctx context.Context, hostID int64) (int, error) {
	db := postgres.GetTx(ctx, r.db)

	today := time.Now().Truncate(24 * time.Hour)
	var count int64
	if err := db.Model(&model.VisitRecord{}).
		Where("host_id = ? AND visited_at >= ?", hostID, today).
		Count(&count).Error; err != nil {
		return 0, err
	}

	return int(count), nil
}

// HasVisitedToday 今天是否已拜访过
func (r *VisitRepository) HasVisitedToday(ctx context.Context, visitorID, hostID int64) (bool, error) {
	db := postgres.GetTx(ctx, r.db)

	today := time.Now().Truncate(24 * time.Hour)
	var count int64
	if err := db.Model(&model.VisitRecord{}).
		Where("visitor_id = ? AND host_id = ? AND visited_at >= ?", visitorID, hostID, today).
		Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

