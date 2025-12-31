// Package handler HTTP 处理器
package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"pets-server/internal/interfaces/http/middleware"
	"pets-server/internal/pkg/response"
)

// ItemHandler 道具处理器
type ItemHandler struct {
	// TODO: 注入 item 应用服务
}

// NewItemHandler 创建道具处理器
func NewItemHandler() *ItemHandler {
	return &ItemHandler{}
}

// RegisterRoutes 注册路由
func (h *ItemHandler) RegisterRoutes(r *gin.RouterGroup) {
	r.GET("", h.GetItems)           // 获取背包
	r.POST("/use", h.UseItem)       // 使用道具
	r.GET("/shop", h.GetShopItems)  // 获取商店道具
	r.POST("/buy", h.BuyItem)       // 购买道具
}

// GetItems 获取背包道具
// GET /api/items
func (h *ItemHandler) GetItems(c *gin.Context) {
	userID := middleware.GetUserID(c)
	_ = userID

	// TODO: 实现获取背包逻辑
	response.Success(c, gin.H{"items": []interface{}{}})
}

// UseItem 使用道具
// POST /api/items/use
func (h *ItemHandler) UseItem(c *gin.Context) {
	// TODO: 实现使用道具逻辑
	response.Success(c, gin.H{"message": "道具使用成功"})
}

// GetShopItems 获取商店道具
// GET /api/items/shop
func (h *ItemHandler) GetShopItems(c *gin.Context) {
	// TODO: 实现获取商店道具逻辑
	response.Success(c, gin.H{"items": []interface{}{}})
}

// BuyItem 购买道具
// POST /api/items/buy
func (h *ItemHandler) BuyItem(c *gin.Context) {
	type BuyRequest struct {
		ItemID   int `json:"itemId" binding:"required"`
		Quantity int `json:"quantity" binding:"required,min=1"`
	}

	var req BuyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	// TODO: 实现购买道具逻辑
	response.Success(c, gin.H{"message": "购买成功"})
}

