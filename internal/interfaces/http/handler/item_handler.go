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
// @Summary      获取背包道具
// @Description  获取当前用户背包中的所有道具列表
// @Tags         item
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Success      200 {object} response.Response{data=object} "获取成功"
// @Failure      500 {object} response.Response "服务器错误"
// @Router       /items [get]
func (h *ItemHandler) GetItems(c *gin.Context) {
	userID := middleware.GetUserID(c)
	_ = userID

	// TODO: 实现获取背包逻辑
	response.Success(c, gin.H{"items": []interface{}{}})
}

// UseItem 使用道具
// @Summary      使用道具
// @Description  使用背包中的道具
// @Tags         item
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Success      200 {object} response.Response "使用成功"
// @Failure      400 {object} response.Response "请求参数错误"
// @Failure      500 {object} response.Response "服务器错误"
// @Router       /items/use [post]
func (h *ItemHandler) UseItem(c *gin.Context) {
	// TODO: 实现使用道具逻辑
	response.Success(c, gin.H{"message": "道具使用成功"})
}

// GetShopItems 获取商店道具
// @Summary      获取商店道具
// @Description  获取商店中可购买的道具列表
// @Tags         item
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Success      200 {object} response.Response{data=object} "获取成功"
// @Failure      500 {object} response.Response "服务器错误"
// @Router       /items/shop [get]
func (h *ItemHandler) GetShopItems(c *gin.Context) {
	// TODO: 实现获取商店道具逻辑
	response.Success(c, gin.H{"items": []interface{}{}})
}

// BuyRequest 购买道具请求
type BuyRequest struct {
	ItemID   int `json:"itemId" binding:"required"`   // 道具ID
	Quantity int `json:"quantity" binding:"required,min=1"` // 购买数量
}

// BuyItem 购买道具
// @Summary      购买道具
// @Description  在商店购买道具
// @Tags         item
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        request body BuyRequest true "购买请求"
// @Success      200 {object} response.Response "购买成功"
// @Failure      400 {object} response.Response "请求参数错误"
// @Failure      500 {object} response.Response "服务器错误"
// @Router       /items/buy [post]
func (h *ItemHandler) BuyItem(c *gin.Context) {

	var req BuyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	// TODO: 实现购买道具逻辑
	response.Success(c, gin.H{"message": "购买成功"})
}

