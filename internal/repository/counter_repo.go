package repository

import (
	"uniS/internal/model"

	"gorm.io/gorm"
)

type CounterRepository interface {
	// Insert 插入一条新记录，count 为本时间段的访问增量。
	Insert(url string, count int64) error
	// Total 查询指定 url 的历史总访问量（SUM）。
	Total(url string) (int64, error)
}

type counterRepository struct {
	db *gorm.DB
}

func NewCounterRepository(db *gorm.DB) CounterRepository {
	return &counterRepository{db: db}
}

// Insert 每 10 分钟调用一次，将增量作为新行写入。
func (r *counterRepository) Insert(url string, count int64) error {
	return r.db.Create(&model.Counter{URL: url, Count: count}).Error
}

// Total 返回该 url 所有记录的 count 总和。
func (r *counterRepository) Total(url string) (int64, error) {
	var total int64
	err := r.db.Model(&model.Counter{}).
		Where("url = ?", url).
		Select("COALESCE(SUM(count), 0)").
		Scan(&total).Error
	return total, err
}
