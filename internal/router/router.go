package router

import (
	"uniS/internal/handler"
	"uniS/internal/middleware"

	"github.com/gin-gonic/gin"
)

func Setup(authHandler *handler.AuthHandler, counterHandler *handler.CounterHandler, jwtSecret string) *gin.Engine {
	r := gin.Default()

	// 公开路由（无需登录）
	auth := r.Group("/auth")
	{
		auth.POST("/wx-login", authHandler.WxLogin)
	}

	// 游客计数器
	r.GET("/counter", counterHandler.GetCounter)

	// 需要鉴权的路由
	api := r.Group("/api", middleware.JWTAuth(jwtSecret))
	{
		api.GET("/me", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"user_id": c.GetUint("user_id"),
				"open_id": c.GetString("open_id"),
			})
		})
	}

	return r
}
