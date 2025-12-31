// Package handler HTTP 处理器
// 这是完整的读写示例
package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	petApp "pets-server/internal/application/pet"
	"pets-server/internal/interfaces/http/middleware"
	"pets-server/internal/pkg/response"
)

// PetHandler 宠物处理器
type PetHandler struct {
	petService *petApp.Service
}

// NewPetHandler 创建宠物处理器
func NewPetHandler(petService *petApp.Service) *PetHandler {
	return &PetHandler{petService: petService}
}

// RegisterRoutes 注册路由
func (h *PetHandler) RegisterRoutes(r *gin.RouterGroup) {
	r.GET("", h.GetMyPet)         // 获取我的宠物（读操作示例）
	r.POST("", h.CreatePet)       // 创建宠物
	r.POST("/feed", h.Feed)       // 喂食（写操作示例）
	r.POST("/play", h.Play)       // 玩耍
	r.POST("/clean", h.Clean)     // 清洁
}

// ============================================================
// 读操作示例：获取宠物详情 GET /api/pet
// 调用链路：
//   Handler.GetMyPet()
//     → AppService.GetPetDetail()
//       → Redis.Get() 尝试读缓存
//       → 缓存未命中
//         → PetRepo.FindByUserID() 查询数据库
//         → 组装 DTO
//       → Redis.Set() 写入缓存
//     ← 返回 DTO
// ============================================================

// GetMyPet 获取我的宠物
// GET /api/pet
func (h *PetHandler) GetMyPet(c *gin.Context) {
	// 1. 从中间件获取用户ID
	userID := middleware.GetUserID(c)

	// 2. 调用应用服务获取宠物详情
	// 应用服务内部会：
	//   - 先查 Redis 缓存
	//   - 缓存未命中则查数据库
	//   - 将结果写入缓存
	pet, err := h.petService.GetPetDetail(c.Request.Context(), userID)
	if err != nil {
		if err == petApp.ErrPetNotFound {
			response.Error(c, http.StatusNotFound, "还没有宠物，快去领养一只吧！")
			return
		}
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	// 3. 返回响应
	response.Success(c, pet)
}

// ============================================================
// 写操作示例：喂食宠物 POST /api/pet/feed
// 调用链路：
//   Handler.Feed()
//     → AppService.FeedPet()
//       → UoW.Do() 开启事务
//         → PetRepo.FindByUserID() 获取宠物
//         → ItemRepo.FindByUserAndItem() 获取食物
//         → Item.Consume() 领域逻辑：扣道具
//         → Pet.Feed() 领域逻辑：喂食
//         → ItemRepo.Save() 保存道具
//         → PetRepo.Save() 保存宠物
//       → 事务提交
//       → Cache.Delete() 清除缓存
//       → EventPublisher.Publish(PetFedEvent) 发布事件
//     ← 返回结果
// ============================================================

// Feed 喂食宠物
// POST /api/pet/feed
func (h *PetHandler) Feed(c *gin.Context) {
	// 1. 获取用户ID
	userID := middleware.GetUserID(c)

	// 2. 绑定请求参数
	var req petApp.FeedPetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	// 3. 调用应用服务执行喂食
	// 应用服务内部会：
	//   - 开启数据库事务
	//   - 获取宠物和道具
	//   - 执行领域逻辑（扣道具、喂食）
	//   - 保存变更
	//   - 提交事务
	//   - 清除 Redis 缓存
	//   - 发布领域事件到 MQ
	result, err := h.petService.FeedPet(c.Request.Context(), userID, req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	// 4. 返回响应
	response.Success(c, result)
}

// CreatePet 创建宠物
// POST /api/pet
func (h *PetHandler) CreatePet(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req petApp.CreatePetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.petService.CreatePet(c.Request.Context(), userID, req)
	if err != nil {
		if err == petApp.ErrAlreadyHasPet {
			response.Error(c, http.StatusConflict, err.Error())
			return
		}
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, result)
}

// Play 和宠物玩耍
// POST /api/pet/play
func (h *PetHandler) Play(c *gin.Context) {
	userID := middleware.GetUserID(c)

	result, err := h.petService.PlayWithPet(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c, result)
}

// Clean 清洁宠物
// POST /api/pet/clean
func (h *PetHandler) Clean(c *gin.Context) {
	userID := middleware.GetUserID(c)

	result, err := h.petService.CleanPet(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c, result)
}

