package repository

import (
	"uniS/internal/model"

	"gorm.io/gorm"
)

type UserRepository interface {
	FindByOpenID(openID string) (*model.User, error)
	Upsert(user *model.User) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) FindByOpenID(openID string) (*model.User, error) {
	var user model.User
	err := r.db.Where("open_id = ?", openID).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Upsert 存在则更新 session_key，不存在则创建
func (r *userRepository) Upsert(user *model.User) error {
	return r.db.Where(model.User{OpenID: user.OpenID}).
		Assign(model.User{
			SessionKey: user.SessionKey,
			UnionID:    user.UnionID,
			NickName:   user.NickName,
			AvatarURL:  user.AvatarURL,
			Gender:     user.Gender,
			Country:    user.Country,
			Province:   user.Province,
			City:       user.City,
		}).
		FirstOrCreate(user).Error
}
