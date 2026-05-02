package repository

import (
	"uniS/internal/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type CounterRepository interface {
	// Get 查询指定 url 的当前计数，不存在返回 count=0。
	Get(url string) (*model.Counter, error)
	// Save 将内存中的累计值直接写入数据库（UPSERT），用于定时落库。
	Save(url string, count int64) error
}

type counterRepository struct {
	db *gorm.DB
}

func NewCounterRepository(db *gorm.DB) CounterRepository {
	return &counterRepository{db: db}
}

func (r *counterRepository) Get(url string) (*model.Counter, error) {
	var c model.Counter
	err := r.db.Where("url = ?", url).First(&c).Error
	if err == gorm.ErrRecordNotFound {
		return &model.Counter{URL: url, Count: 0}, nil
	}
	if err != nil {
		return nil, err
	}
	return &c, nil
}

// Save UPSERT：存在则更新 count，不存在则插入。
// count 直接使用内存中的累计值，不做 +1 运算。
func (r *counterRepository) Save(url string, count int64) error {
	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "url"}},
		DoUpdates: clause.AssignmentColumns([]string{"count", "updated_at"}),
	}).Create(&model.Counter{URL: url, Count: count}).Error
}
