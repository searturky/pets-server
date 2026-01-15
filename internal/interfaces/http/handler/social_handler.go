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
// @Summary      获取好友列表
// @Description  获取当前用户的好友列表，包括亲密度等信息
// @Tags         social
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Success      200 {object} response.Response{data=social.FriendListResponse} "获取成功"
// @Failure      500 {object} response.Response "服务器错误"
// @Router       /social/friends [get]
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
// @Summary      发送好友申请
// @Description  向指定用户发送好友申请
// @Tags         social
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        request body social.AddFriendRequest true "好友申请"
// @Success      200 {object} response.Response "申请已发送"
// @Failure      400 {object} response.Response "请求参数错误"
// @Failure      500 {object} response.Response "服务器错误"
// @Router       /social/friends/request [post]
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
// @Summary      接受好友申请
// @Description  接受指定的好友申请
// @Tags         social
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        id path int true "好友申请ID"
// @Success      200 {object} response.Response "已添加为好友"
// @Failure      400 {object} response.Response "请求参数错误"
// @Failure      500 {object} response.Response "服务器错误"
// @Router       /social/friends/accept/{id} [post]
func (h *SocialHandler) AcceptFriendRequest(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var friendshipID int
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
// @Summary      获取好友申请列表
// @Description  获取待处理的好友申请列表
// @Tags         social
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Success      200 {object} response.Response{data=object} "获取成功"
// @Failure      500 {object} response.Response "服务器错误"
// @Router       /social/friends/requests [get]
func (h *SocialHandler) GetFriendRequests(c *gin.Context) {
	// TODO: 实现获取好友申请列表
	response.Success(c, gin.H{"requests": []interface{}{}})
}

// --- 礼物相关 ---

// GetReceivedGifts 获取收到的礼物
// @Summary      获取收到的礼物
// @Description  获取当前用户收到的礼物列表
// @Tags         social
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Success      200 {object} response.Response{data=object} "获取成功"
// @Failure      500 {object} response.Response "服务器错误"
// @Router       /social/gifts [get]
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
// @Summary      发送礼物
// @Description  向好友发送礼物
// @Tags         social
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        request body social.SendGiftRequest true "礼物请求"
// @Success      200 {object} response.Response "发送成功"
// @Failure      400 {object} response.Response "请求参数错误"
// @Failure      500 {object} response.Response "服务器错误"
// @Router       /social/gifts [post]
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
// @Summary      拜访好友
// @Description  拜访好友的宠物，可获得奖励
// @Tags         social
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        userId path int true "好友用户ID"
// @Success      200 {object} response.Response "拜访成功"
// @Failure      400 {object} response.Response "请求参数错误"
// @Failure      500 {object} response.Response "服务器错误"
// @Router       /social/visit/{userId} [get]
func (h *SocialHandler) VisitFriend(c *gin.Context) {
	// TODO: 实现拜访好友逻辑
	response.Success(c, gin.H{"message": "拜访成功"})
}

// --- 交易相关 ---

// GetTrades 获取交易列表
// @Summary      获取交易列表
// @Description  获取当前用户的交易记录
// @Tags         social
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Success      200 {object} response.Response{data=object} "获取成功"
// @Failure      500 {object} response.Response "服务器错误"
// @Router       /social/trades [get]
func (h *SocialHandler) GetTrades(c *gin.Context) {
	// TODO: 实现获取交易列表
	response.Success(c, gin.H{"trades": []interface{}{}})
}

// CreateTrade 创建交易
// @Summary      创建交易
// @Description  创建一个新的道具交易请求
// @Tags         social
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        request body social.CreateTradeRequest true "交易请求"
// @Success      200 {object} response.Response "创建成功"
// @Failure      400 {object} response.Response "请求参数错误"
// @Failure      500 {object} response.Response "服务器错误"
// @Router       /social/trades [post]
func (h *SocialHandler) CreateTrade(c *gin.Context) {
	// TODO: 实现创建交易
	response.Success(c, gin.H{"message": "交易已创建"})
}

// AcceptTrade 接受交易
// @Summary      接受交易
// @Description  接受指定的交易请求
// @Tags         social
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        id path int true "交易ID"
// @Success      200 {object} response.Response "交易完成"
// @Failure      400 {object} response.Response "请求参数错误"
// @Failure      500 {object} response.Response "服务器错误"
// @Router       /social/trades/accept/{id} [post]
func (h *SocialHandler) AcceptTrade(c *gin.Context) {
	// TODO: 实现接受交易
	response.Success(c, gin.H{"message": "交易完成"})
}

// CancelTrade 取消交易
// @Summary      取消交易
// @Description  取消指定的交易请求
// @Tags         social
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        id path int true "交易ID"
// @Success      200 {object} response.Response "交易已取消"
// @Failure      400 {object} response.Response "请求参数错误"
// @Failure      500 {object} response.Response "服务器错误"
// @Router       /social/trades/cancel/{id} [post]
func (h *SocialHandler) CancelTrade(c *gin.Context) {
	// TODO: 实现取消交易
	response.Success(c, gin.H{"message": "交易已取消"})
}

