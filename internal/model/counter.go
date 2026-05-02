package model

import "time"

// Counter 游客访问计数器
// 每 10 分钟插入一条新记录，记录该时间段内的访问增量。
// 总访问量 = SELECT SUM(count) FROM counters WHERE url = ?
type Counter struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	CreatedAt time.Time `                                json:"created_at"`
	UpdatedAt time.Time `                                json:"updated_at"`
	URL       string    `gorm:"index;size:256;not null"  json:"url"`   // 普通索引，允许多行
	Count     int64     `gorm:"default:0;not null"       json:"count"` // 本段增量
}
