package model

import "gorm.io/gorm"

// User 微信用户表
type User struct {
	gorm.Model
	OpenID     string `gorm:"uniqueIndex;size:64;not null" json:"open_id"`
	UnionID    string `gorm:"index;size:64"               json:"union_id"`
	NickName   string `gorm:"size:64"                     json:"nick_name"`
	AvatarURL  string `gorm:"size:512"                    json:"avatar_url"`
	Gender     int    `gorm:"default:0"                   json:"gender"` // 0未知 1男 2女
	Country    string `gorm:"size:64"                     json:"country"`
	Province   string `gorm:"size:64"                     json:"province"`
	City       string `gorm:"size:64"                     json:"city"`
	SessionKey string `gorm:"size:128"                    json:"-"` // 不对外暴露
}
