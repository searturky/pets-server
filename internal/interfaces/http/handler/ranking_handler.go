// Package handler HTTP 处理器
package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"pets-server/internal/application/ranking"
	"pets-server/internal/interfaces/http/middleware"
	"pets-server/internal/pkg/response"
)

// RankingHandler 排行榜处理器
type RankingHandler struct {
	rankingService *ranking.Service
}

// NewRankingHandler 创建排行榜处理器
func NewRankingHandler(rankingService *ranking.Service) *RankingHandler {
	return &RankingHandler{rankingService: rankingService}
}

// RegisterRoutes 注册路由
func (h *RankingHandler) RegisterRoutes(r *gin.RouterGroup) {
	r.GET("", h.GetRanking)
}

// GetRanking 获取排行榜
// GET /api/ranking?type=pet_level&offset=0&limit=20
func (h *RankingHandler) GetRanking(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req ranking.RankingRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.rankingService.GetRanking(c.Request.Context(), userID, req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, result)
}

