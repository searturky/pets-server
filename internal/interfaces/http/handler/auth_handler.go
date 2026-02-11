// Package handler HTTP 处理器
package handler

import (
	"github.com/gin-gonic/gin"

	"pets-server/internal/application/auth"
	"pets-server/internal/interfaces/http/middleware"
	"pets-server/internal/pkg/response"
)

// AuthHandler 认证处理器
type AuthHandler struct {
	authService *auth.Service
}

// NewAuthHandler 创建认证处理器
func NewAuthHandler(authService *auth.Service) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// RegisterRoutes 注册路由
func (h *AuthHandler) RegisterRoutes(r *gin.RouterGroup) {
	r.POST("/wx-login", h.WxLogin)
	r.POST("/login", h.Login)
	r.POST("/register", h.Register)
}

// RegisterAuthRoutes 注册需要认证的路由
func (h *AuthHandler) RegisterAuthRoutes(r *gin.RouterGroup) {
	r.GET("/me", h.GetCurrentUser)
	r.POST("/logout", h.Logout)
	r.POST("/kick", h.KickUser)
}

// Login 账号密码登录
// @Summary      账号密码登录
// @Description  使用用户名和密码登录获取 JWT token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body auth.LoginRequest true "登录请求"
// @Success      200 {object} response.Response{data=auth.LoginResponse} "登录成功"
// @Failure      400 {object} response.Response "请求参数错误"
// @Failure      401 {object} response.Response "用户名或密码错误"
// @Failure      500 {object} response.Response "服务器错误"
// @Router       /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req auth.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.CodeBadRequest, err.Error())
		return
	}

	result, err := h.authService.Login(c.Request.Context(), req)
	if err != nil {
		response.Error(c, response.CodeBadRequest, err.Error())
		return
	}

	response.Success(c, result)
}

// Register 账号注册
// @Summary      用户注册
// @Description  注册新用户账号
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body auth.RegisterRequest true "注册请求"
// @Success      200 {object} response.Response{data=auth.LoginResponse} "注册成功"
// @Failure      400 {object} response.Response "请求参数错误"
// @Failure      409 {object} response.Response "用户名已存在"
// @Failure      500 {object} response.Response "服务器错误"
// @Router       /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req auth.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.CodeBadRequest, err.Error())
		return
	}

	result, err := h.authService.Register(c.Request.Context(), req)
	if err != nil {
		response.Error(c, response.CodeBadRequest, err.Error())
		return
	}

	response.Success(c, result)
}

// WxLogin 微信登录
// @Summary      微信登录
// @Description  通过微信 code 登录获取 JWT token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body auth.WxLoginRequest true "微信登录请求"
// @Success      200 {object} response.Response{data=auth.WxLoginResponse} "登录成功"
// @Failure      400 {object} response.Response "请求参数错误"
// @Failure      500 {object} response.Response "服务器错误"
// @Router       /auth/wx-login [post]
func (h *AuthHandler) WxLogin(c *gin.Context) {
	var req auth.WxLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.CodeBadRequest, err.Error())
		return
	}

	result, err := h.authService.WxLogin(c.Request.Context(), req)
	if err != nil {
		response.Error(c, response.CodeInternalError, err.Error())
		return
	}

	response.Success(c, result)
}

// GetCurrentUser 获取当前用户信息
// @Summary      获取当前用户信息
// @Description  获取当前登录用户的详细信息
// @Tags         auth
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Success      200 {object} response.Response{data=auth.UserInfo} "获取成功"
// @Failure      401 {object} response.Response "未授权"
// @Failure      500 {object} response.Response "服务器错误"
// @Router       /auth/me [get]
func (h *AuthHandler) GetCurrentUser(c *gin.Context) {
	userID := middleware.GetUserID(c)

	userInfo, err := h.authService.GetUserInfo(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, response.CodeInternalError, err.Error())
		return
	}

	response.Success(c, userInfo)
}

// Logout 退出登录
// @Summary      退出登录
// @Description  清除当前用户会话
// @Tags         auth
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Success      200 {object} response.Response "退出成功"
// @Failure      401 {object} response.Response "未授权"
// @Failure      500 {object} response.Response "服务器错误"
// @Router       /auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "unauthorized")
		return
	}

	if err := h.authService.Logout(c.Request.Context(), userID); err != nil {
		response.Error(c, response.CodeInternalError, err.Error())
		return
	}

	response.SuccessWithMessage(c, "logout success", nil)
}

// KickUser 主动踢下线
// @Summary      主动踢下线
// @Description  清除指定用户会话
// @Tags         auth
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        request body auth.KickRequest true "踢下线请求"
// @Success      200 {object} response.Response "踢下线成功"
// @Failure      400 {object} response.Response "请求参数错误"
// @Failure      401 {object} response.Response "未授权"
// @Failure      500 {object} response.Response "服务器错误"
// @Router       /auth/kick [post]
func (h *AuthHandler) KickUser(c *gin.Context) {
	var req auth.KickRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.CodeBadRequest, err.Error())
		return
	}
	if req.UserID <= 0 {
		response.Error(c, response.CodeBadRequest, "invalid user id")
		return
	}

	if err := h.authService.KickUser(c.Request.Context(), req.UserID); err != nil {
		response.Error(c, response.CodeInternalError, err.Error())
		return
	}

	response.SuccessWithMessage(c, "kick success", nil)
}
