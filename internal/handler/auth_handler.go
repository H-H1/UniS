package handler

import (
	"net/http"
	"uniS/internal/service"
	"uniS/pkg/logger"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authSvc service.AuthService
}

func NewAuthHandler(authSvc service.AuthService) *AuthHandler {
	return &AuthHandler{authSvc: authSvc}
}

// WxLogin godoc
// POST /auth/wx-login
func (h *AuthHandler) WxLogin(c *gin.Context) {
	var req service.WxLoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("auth_handler", "请求参数绑定失败", map[string]any{
			"client_ip": c.ClientIP(),
			"error":     err.Error(),
		})
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": err.Error()})
		return
	}

	logger.Info("auth_handler", "收到 WxLogin 请求", map[string]any{
		"client_ip": c.ClientIP(),
		"nick_name": req.NickName,
		"country":   req.Country,
		"province":  req.Province,
		"city":      req.City,
	})

	resp, err := h.authSvc.WxLogin(&req)
	if err != nil {
		logger.Error("auth_handler", "WxLogin 处理失败", map[string]any{
			"client_ip": c.ClientIP(),
			"nick_name": req.NickName,
			"error":     err.Error(),
		})
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": err.Error()})
		return
	}

	logger.Info("auth_handler", "WxLogin 响应成功", map[string]any{
		"client_ip": c.ClientIP(),
		"user_id":   resp.User.ID,
		"openid":    resp.User.OpenID,
		"nick_name": resp.User.NickName,
	})

	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "ok", "data": resp})
}
