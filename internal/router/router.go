package router

import (
	"uniS/internal/handler"
	"uniS/internal/middleware"

	"github.com/gin-gonic/gin"
)

func Setup(authHandler *handler.AuthHandler, counterHandler *handler.CounterHandler, userHandler *handler.UserHandler, testHandler *handler.TestHandler, jwtSecret string) *gin.Engine {
	r := gin.Default()

	// 公开路由（无需登录）
	auth := r.Group("/auth")
	{
		auth.POST("/wx-login", authHandler.WxLogin)
		auth.POST("/profile", authHandler.WxLogin)
	}

	// 游客计数器
	r.GET("/counter", counterHandler.GetCounter)

	// 需要鉴权的路由
	api := r.Group("/", middleware.JWTAuth(jwtSecret))
	{
		api.GET("/me", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"user_id": c.GetUint("user_id"),
				"open_id": c.GetString("open_id"),
			})
		})

		// 用户信息
		user := api.Group("/user")
		{
			user.POST("/profile", userHandler.UpdateProfile) // 更新昵称/头像URL
			user.POST("/avatar", userHandler.UploadAvatar)   // 上传头像文件
		}

		// 测试相关
		test := api.Group("/test")
		{
			test.POST("/submit", testHandler.SubmitTest)
			test.GET("/history", testHandler.GetTestHistory)
			test.GET("/:id", testHandler.GetTestResult)
		}
	}

	// 静态文件服务（头像访问）
	r.Static("/static/avatars", "./uploads/avatars")

	return r
}
