package repository

import (
	"uniS/internal/model"

	"gorm.io/gorm"
)

type TestResultRepository interface {
	Create(result *model.TestResult) error
	FindByUserID(userID uint) ([]model.TestResult, error)
	FindByID(id uint) (*model.TestResult, error)
}

type testResultRepository struct {
	db *gorm.DB
}

func NewTestResultRepository(db *gorm.DB) TestResultRepository {
	return &testResultRepository{db: db}
}

func (r *testResultRepository) Create(result *model.TestResult) error {
	return r.db.Create(result).Error
}

func (r *testResultRepository) FindByUserID(userID uint) ([]model.TestResult, error) {
	var results []model.TestResult
	err := r.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&results).Error
	return results, err
}

func (r *testResultRepository) FindByID(id uint) (*model.TestResult, error) {
	var result model.TestResult
	err := r.db.First(&result, id).Error
	if err != nil {
		return nil, err
	}
	return &result, nil
}
