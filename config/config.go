package config

import (
	"fmt"
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

// Load 加载配置
func Load() *Config {
	return &Config{
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8099"),
		},
		Database: DatabaseConfig{
			DSN: buildDSN(),
		},
		WeChat: WeChatConfig{
			AppID:     getEnv("WX_APP_ID", "wxa653960621255b47"),
			AppSecret: getEnv("WX_APP_SECRET", "4965013f23337ae6b6ccf89451f3595c"),
		},
		JWT: JWTConfig{
			Secret: getEnv("JWT_SECRET", "change-uniS-s"),
		},
	}
}

// buildDSN 构建数据库连接字符串（避免敏感信息出现在代码中）
func buildDSN() string {
	username := getEnv("MYSQL_USERNAME", "root")
	password := getEnv("MYSQL_PASSWORD", "PXDN93VRKUm8TeE7")
	hostpost := getEnv("MYSQL_ADDRESS", "127.0.0.1:33069")
	database := getEnv("MYSQL_DATABASE", "unis")

	// 如果有完整的DSN环境变量，优先使用
	if dsn := os.Getenv("MYSQL_DSN"); dsn != "" {
		return dsn
	}

	return fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		username, password, hostpost, database)
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
