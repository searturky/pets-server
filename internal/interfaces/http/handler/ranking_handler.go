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
// @Summary      获取排行榜
// @Description  获取排行榜数据，支持宠物等级、成就数量、亲密度等类型
// @Tags         ranking
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        type query string true "排行榜类型" Enums(pet_level, achievement, intimacy)
// @Param        offset query int false "偏移量" default(0)
// @Param        limit query int false "数量限制" default(20)
// @Success      200 {object} response.Response{data=ranking.RankingResponse} "获取成功"
// @Failure      400 {object} response.Response "请求参数错误"
// @Failure      500 {object} response.Response "服务器错误"
// @Router       /ranking [get]
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

