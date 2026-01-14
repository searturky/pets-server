// Package handler HTTP 处理器
package handler

import (
	"net/http"

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
}

// Login 账号密码登录
// POST /api/v1/auth/login
func (h *AuthHandler) Login(c *gin.Context) {
	var req auth.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.authService.Login(c.Request.Context(), req)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, err.Error())
		return
	}

	response.Success(c, result)
}

// Register 账号注册
// POST /api/v1/auth/register
func (h *AuthHandler) Register(c *gin.Context) {
	var req auth.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.authService.Register(c.Request.Context(), req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c, result)
}

// WxLogin 微信登录
// POST /api/v1/auth/wx-login
func (h *AuthHandler) WxLogin(c *gin.Context) {
	var req auth.WxLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.authService.WxLogin(c.Request.Context(), req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, result)
}

// GetCurrentUser 获取当前用户信息
// GET /api/v1/auth/me
func (h *AuthHandler) GetCurrentUser(c *gin.Context) {
	userID := middleware.GetUserID(c)

	userInfo, err := h.authService.GetUserInfo(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, userInfo)
}
