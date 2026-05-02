package service

import (
	"uniS/internal/model"
	"uniS/internal/repository"
	"uniS/pkg/logger"
)

type UpdateProfileReq struct {
	NickName  string `json:"nickName"`
	AvatarURL string `json:"avatarUrl"`
}

type UserService interface {
	// UpdateProfile 更新昵称和头像（仅更新非空字段）
	UpdateProfile(userID uint, req *UpdateProfileReq) (*model.User, error)
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

func (s *userService) UpdateProfile(userID uint, req *UpdateProfileReq) (*model.User, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		logger.Error("user_service", "FindByID 失败", map[string]any{
			"user_id": userID,
			"error":   err.Error(),
		})
		return nil, err
	}

	// 只更新传入的非空字段
	if req.NickName != "" {
		user.NickName = req.NickName
	}
	if req.AvatarURL != "" {
		user.AvatarURL = req.AvatarURL
	}

	if err := s.userRepo.UpdateProfile(user); err != nil {
		logger.Error("user_service", "UpdateProfile 失败", map[string]any{
			"user_id": userID,
			"error":   err.Error(),
		})
		return nil, err
	}

	logger.Info("user_service", "UpdateProfile 成功", map[string]any{
		"user_id":    userID,
		"nick_name":  user.NickName,
		"avatar_url": user.AvatarURL,
	})
	return user, nil
}
