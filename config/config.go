package config

import (
	"os"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	WeChat   WeChatConfig
	JWT      JWTConfig
}

type ServerConfig struct {
	Port string
}

type DatabaseConfig struct {
	DSN string
}

type WeChatConfig struct {
	AppID     string
	AppSecret string
}

type JWTConfig struct {
	Secret string
}

func Load() *Config {
	return &Config{
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8099"),
		},
		Database: DatabaseConfig{
			DSN: getEnv("DB_DSN", "root:PXDN93VRKUm8TeE7@tcp(127.0.0.1:33069)/unis?charset=utf8mb4&parseTime=True&loc=Local"),
		},
		WeChat: WeChatConfig{
			AppID:     getEnv("WX_APP_ID", "wx6445e83e8e9c9885"),
			AppSecret: getEnv("WX_APP_SECRET", "4965013f23337ae6b6ccf89451f3595c"),
		},
		JWT: JWTConfig{
			Secret: getEnv("JWT_SECRET", "change-me-in-production"),
		},
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
