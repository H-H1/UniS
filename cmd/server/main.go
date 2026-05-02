package main

import (
	"log"
	"uniS/config"
	"uniS/internal/database"
	"uniS/internal/handler"
	"uniS/internal/repository"
	"uniS/internal/router"
	"uniS/internal/service"
	"uniS/pkg/wechat"
)

func main() {
	cfg := config.Load()

	// 数据库
	db, err := database.NewDB(cfg.Database.DSN)
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	// 依赖注入（手动 DI，简单清晰）
	wxClient := wechat.NewClient(cfg.WeChat.AppID, cfg.WeChat.AppSecret)
	userRepo := repository.NewUserRepository(db)
	counterRepo := repository.NewCounterRepository(db)
	authSvc := service.NewAuthService(wxClient, userRepo, cfg.JWT.Secret)
	authHandler := handler.NewAuthHandler(authSvc)
	counterHandler := handler.NewCounterHandler(counterRepo)

	// 路由
	r := router.Setup(authHandler, counterHandler, cfg.JWT.Secret)

	log.Printf("server running on :%s", cfg.Server.Port)
	if err := r.Run(":" + cfg.Server.Port); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
