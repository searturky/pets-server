// Package handler HTTP 处理器
// 这是完整的读写示例
package handler

import (
	"strconv"

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
	r.GET("", h.GetMyPet)          // 获取我的宠物（读操作示例）
	r.GET("/status", h.GetStatus)  // 获取轻量状态
	r.PUT("/active", h.SetActive)  // 设置主宠物
	r.POST("", h.CreatePet)        // 创建宠物
	r.POST("/feed", h.Feed)        // 喂食（写操作示例）
	r.POST("/play", h.Play)        // 玩耍
	r.POST("/clean", h.Clean)      // 清洁
}

// ============================================================
// 读操作示例：获取宠物详情 GET /api/pet
// 调用链路：
//   Handler.GetMyPet()
//     → AppService.GetPetDetail()
//       → PetRepo.FindByUserID() 查询数据库
//       → 组装 DTO
//     ← 返回 DTO
// ============================================================

// GetMyPet 获取我的宠物
// @Summary      获取我的宠物
// @Description  获取当前用户的宠物详细信息，包括外观、性格、技能和状态
// @Tags         pet
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Success      200 {object} response.Response{data=petApp.PetDetailDTO} "获取成功"
// @Failure      404 {object} response.Response "还没有宠物"
// @Failure      500 {object} response.Response "服务器错误"
// @Router       /pet [get]
func (h *PetHandler) GetMyPet(c *gin.Context) {
	// 1. 从中间件获取用户ID
	userID := middleware.GetUserID(c)

	// 2. 调用应用服务获取宠物详情
	// 应用服务内部会直接查询数据库并组装响应
	pet, err := h.petService.GetPetDetail(c.Request.Context(), userID)
	if err != nil {
		if err == petApp.ErrPetNotFound {
			response.SuccessWithMessageAndCode(c, response.CodePetNotFound, "还没有宠物，快去领养一只吧！", nil)
			return
		}
		response.Error(c, response.CodeInternalError, err.Error())
		return
	}

	// 3. 返回响应
	response.Success(c, pet)
}

// GetStatus 获取我的宠物轻量状态
// @Summary      获取宠物轻量状态
// @Description  获取当前用户宠物的轻量状态信息（实时状态与版本信息）
// @Tags         pet
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        petId query int true "宠物ID"
// @Success      200 {object} response.Response{data=petApp.PetStatusDTO} "获取成功"
// @Failure      400 {object} response.Response "请求参数错误"
// @Failure      404 {object} response.Response "还没有宠物"
// @Failure      403 {object} response.Response "非宠物主人"
// @Failure      500 {object} response.Response "服务器错误"
// @Router       /pet/status [get]
func (h *PetHandler) GetStatus(c *gin.Context) {
	userID := middleware.GetUserID(c)
	petIDStr := c.Query("petId")
	petID, err := strconv.Atoi(petIDStr)
	if err != nil || petID <= 0 {
		response.Error(c, response.CodeBadRequest, "petId 参数无效")
		return
	}

	status, err := h.petService.GetPetStatus(c.Request.Context(), userID, petID)
	if err != nil {
		if err == petApp.ErrPetNotFound {
			response.SuccessWithMessageAndCode(c, response.CodePetNotFound, "还没有宠物，快去领养一只吧！", nil)
			return
		}
		if err == petApp.ErrNotPetOwner {
			response.Error(c, response.CodeForbidden, err.Error())
			return
		}
		response.Error(c, response.CodeInternalError, err.Error())
		return
	}

	response.Success(c, status)
}

// SetActive 设置主宠物
// @Summary      设置主宠物
// @Description  设置当前用户的主宠物（主页默认展示和互动对象）
// @Tags         pet
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        request body petApp.SetActivePetRequest true "设置主宠物请求"
// @Success      200 {object} response.Response{data=petApp.SetActivePetResponse} "设置成功"
// @Failure      400 {object} response.Response "请求参数错误"
// @Failure      403 {object} response.Response "非宠物主人"
// @Failure      404 {object} response.Response "宠物不存在"
// @Failure      500 {object} response.Response "服务器错误"
// @Router       /pet/active [put]
func (h *PetHandler) SetActive(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req petApp.SetActivePetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.CodeBadRequest, err.Error())
		return
	}

	result, err := h.petService.SetActivePet(c.Request.Context(), userID, req.PetID)
	if err != nil {
		if err == petApp.ErrPetNotFound {
			response.SuccessWithMessageAndCode(c, response.CodePetNotFound, "宠物不存在", nil)
			return
		}
		if err == petApp.ErrNotPetOwner {
			response.Error(c, response.CodeForbidden, err.Error())
			return
		}
		response.Error(c, response.CodeInternalError, err.Error())
		return
	}

	response.Success(c, result)
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
// @Summary      喂食宠物
// @Description  使用食物道具喂食宠物，增加饱食度和经验值，可能升级
// @Tags         pet
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        request body petApp.FeedPetRequest true "喂食请求"
// @Success      200 {object} response.Response{data=petApp.FeedPetResponse} "喂食成功"
// @Failure      400 {object} response.Response "请求参数错误"
// @Failure      500 {object} response.Response "服务器错误"
// @Router       /pet/feed [post]
func (h *PetHandler) Feed(c *gin.Context) {
	// 1. 获取用户ID
	userID := middleware.GetUserID(c)

	// 2. 绑定请求参数
	var req petApp.FeedPetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.CodeBadRequest, err.Error())
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
		response.Error(c, response.CodeBadRequest, err.Error())
		return
	}

	// 4. 返回响应
	response.Success(c, result)
}

// CreatePet 创建宠物
// @Summary      创建宠物
// @Description  创建一个新宠物，可以指定物种或随机生成
// @Tags         pet
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        request body petApp.CreatePetRequest true "创建宠物请求"
// @Success      200 {object} response.Response{data=petApp.CreatePetResponse} "创建成功"
// @Failure      400 {object} response.Response "请求参数错误"
// @Failure      409 {object} response.Response "已拥有宠物"
// @Failure      500 {object} response.Response "服务器错误"
// @Router       /pet [post]
func (h *PetHandler) CreatePet(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req petApp.CreatePetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.CodeBadRequest, err.Error())
		return
	}

	result, err := h.petService.CreatePet(c.Request.Context(), userID, req)
	if err != nil {
		if err == petApp.ErrAlreadyHasPet {
			response.Error(c, response.CodeConflict, err.Error())
			return
		}
		response.Error(c, response.CodeInternalError, err.Error())
		return
	}

	response.Success(c, result)
}

// Play 和宠物玩耍
// @Summary      和宠物玩耍
// @Description  与宠物互动玩耍，增加快乐度和经验值，消耗精力
// @Tags         pet
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Success      200 {object} response.Response{data=petApp.PlayPetResponse} "玩耍成功"
// @Failure      400 {object} response.Response "请求参数错误"
// @Failure      500 {object} response.Response "服务器错误"
// @Router       /pet/play [post]
func (h *PetHandler) Play(c *gin.Context) {
	userID := middleware.GetUserID(c)

	result, err := h.petService.PlayWithPet(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, response.CodeBadRequest, err.Error())
		return
	}

	response.Success(c, result)
}

// Clean 清洁宠物
// @Summary      清洁宠物
// @Description  清洁宠物，增加清洁度
// @Tags         pet
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Success      200 {object} response.Response{data=petApp.CleanPetResponse} "清洁成功"
// @Failure      400 {object} response.Response "请求参数错误"
// @Failure      500 {object} response.Response "服务器错误"
// @Router       /pet/clean [post]
func (h *PetHandler) Clean(c *gin.Context) {
	userID := middleware.GetUserID(c)

	result, err := h.petService.CleanPet(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, response.CodeBadRequest, err.Error())
		return
	}

	response.Success(c, result)
}
