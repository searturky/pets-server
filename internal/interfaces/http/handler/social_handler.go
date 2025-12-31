// Package handler HTTP 处理器
package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"pets-server/internal/application/social"
	"pets-server/internal/interfaces/http/middleware"
	"pets-server/internal/pkg/response"
)

// SocialHandler 社交处理器
type SocialHandler struct {
	socialService *social.Service
}

// NewSocialHandler 创建社交处理器
func NewSocialHandler(socialService *social.Service) *SocialHandler {
	return &SocialHandler{socialService: socialService}
}

// RegisterRoutes 注册路由
func (h *SocialHandler) RegisterRoutes(r *gin.RouterGroup) {
	// 好友
	friends := r.Group("/friends")
	{
		friends.GET("", h.GetFriendList)
		friends.POST("/request", h.SendFriendRequest)
		friends.POST("/accept/:id", h.AcceptFriendRequest)
		friends.GET("/requests", h.GetFriendRequests)
	}

	// 礼物
	gifts := r.Group("/gifts")
	{
		gifts.GET("", h.GetReceivedGifts)
		gifts.POST("", h.SendGift)
	}

	// 拜访
	r.GET("/visit/:userId", h.VisitFriend)

	// 交易
	trades := r.Group("/trades")
	{
		trades.GET("", h.GetTrades)
		trades.POST("", h.CreateTrade)
		trades.POST("/accept/:id", h.AcceptTrade)
		trades.POST("/cancel/:id", h.CancelTrade)
	}
}

// --- 好友相关 ---

// GetFriendList 获取好友列表
func (h *SocialHandler) GetFriendList(c *gin.Context) {
	userID := middleware.GetUserID(c)

	result, err := h.socialService.GetFriendList(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, result)
}

// SendFriendRequest 发送好友申请
func (h *SocialHandler) SendFriendRequest(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req social.AddFriendRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.socialService.SendFriendRequest(c.Request.Context(), userID, req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c, gin.H{"message": "好友申请已发送"})
}

// AcceptFriendRequest 接受好友申请
func (h *SocialHandler) AcceptFriendRequest(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var friendshipID int64
	if _, err := c.Params.Get("id"); err {
		response.Error(c, http.StatusBadRequest, "invalid friendship id")
		return
	}

	if err := h.socialService.AcceptFriendRequest(c.Request.Context(), userID, friendshipID); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c, gin.H{"message": "已添加为好友"})
}

// GetFriendRequests 获取待处理的好友申请
func (h *SocialHandler) GetFriendRequests(c *gin.Context) {
	// TODO: 实现获取好友申请列表
	response.Success(c, gin.H{"requests": []interface{}{}})
}

// --- 礼物相关 ---

// GetReceivedGifts 获取收到的礼物
func (h *SocialHandler) GetReceivedGifts(c *gin.Context) {
	userID := middleware.GetUserID(c)

	gifts, err := h.socialService.GetReceivedGifts(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, gin.H{"gifts": gifts})
}

// SendGift 发送礼物
func (h *SocialHandler) SendGift(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req social.SendGiftRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.socialService.SendGift(c.Request.Context(), userID, req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c, gin.H{"message": "礼物发送成功"})
}

// --- 拜访相关 ---

// VisitFriend 拜访好友
func (h *SocialHandler) VisitFriend(c *gin.Context) {
	// TODO: 实现拜访好友逻辑
	response.Success(c, gin.H{"message": "拜访成功"})
}

// --- 交易相关 ---

// GetTrades 获取交易列表
func (h *SocialHandler) GetTrades(c *gin.Context) {
	// TODO: 实现获取交易列表
	response.Success(c, gin.H{"trades": []interface{}{}})
}

// CreateTrade 创建交易
func (h *SocialHandler) CreateTrade(c *gin.Context) {
	// TODO: 实现创建交易
	response.Success(c, gin.H{"message": "交易已创建"})
}

// AcceptTrade 接受交易
func (h *SocialHandler) AcceptTrade(c *gin.Context) {
	// TODO: 实现接受交易
	response.Success(c, gin.H{"message": "交易完成"})
}

// CancelTrade 取消交易
func (h *SocialHandler) CancelTrade(c *gin.Context) {
	// TODO: 实现取消交易
	response.Success(c, gin.H{"message": "交易已取消"})
}

