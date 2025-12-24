package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kronos/backend/internal/service"
)

type AuthHandler struct {
	authService  service.AuthService
	auditService service.AuditService
}

func NewAuthHandler(authService service.AuthService, auditService service.AuditService) *AuthHandler {
	return &AuthHandler{
		authService:  authService,
		auditService: auditService,
	}
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RegisterRequest struct {
	Username    string `json:"username" binding:"required"`
	Email       string `json:"email" binding:"required,email"`
	Password    string `json:"password" binding:"required,min=6"`
	DisplayName string `json:"display_name"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ip := c.ClientIP()
	userAgent := c.Request.UserAgent()
	token, user, err := h.authService.Login(req.Username, req.Password, ip)
	
	// 记录登录审计日志
	status := "success"
	errorMsg := ""
	if err != nil {
		status = "failed"
		errorMsg = err.Error()
	}
	
	if h.auditService != nil {
		userID := uint64(0)
		if user != nil {
			userID = user.ID
		}
		_ = h.auditService.LogOperation(
			&userID,
			req.Username,
			"login",
			"auth",
			nil,
			"",
			ip,
			userAgent,
			nil,
			nil,
			nil,
			nil,
			status,
			errorMsg,
			0,
		)
	}
	
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user": gin.H{
			"id":           user.ID,
			"username":     user.Username,
			"email":        user.Email,
			"display_name": user.DisplayName,
			"roles":        user.Roles,
		},
	})
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.authService.Register(req.Username, req.Email, req.Password, req.DisplayName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"user": gin.H{
			"id":           user.ID,
			"username":     user.Username,
			"email":        user.Email,
			"display_name": user.DisplayName,
		},
	})
}
