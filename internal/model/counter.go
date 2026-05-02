package model

import "time"

// Counter 游客访问计数器
// 每个 URL 对应一行，count 记录累计访问次数。
type Counter struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"  json:"id"`
	CreatedAt time.Time `                                 json:"created_at"`
	UpdatedAt time.Time `                                 json:"updated_at"`
	URL       string    `gorm:"uniqueIndex;size:256;not null" json:"url"`
	Count     int64     `gorm:"default:0;not null"        json:"count"`
}
