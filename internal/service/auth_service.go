package service

import (
	"uniS/internal/model"
	"uniS/internal/repository"
	"uniS/pkg/jwt"
	"uniS/pkg/logger"
	"uniS/pkg/wechat"
)

type WxLoginReq struct {
	Code      string `json:"code"       binding:"required"`
	NickName  string `json:"nickName"`
	AvatarURL string `json:"avatarUrl"`
	Gender    int    `json:"gender"`
	Country   string `json:"country"`
	Province  string `json:"province"`
	City      string `json:"city"`
}

type WxLoginResp struct {
	Token string      `json:"token"`
	User  *model.User `json:"user"`
}

type AuthService interface {
	WxLogin(req *WxLoginReq) (*WxLoginResp, error)
}

type authService struct {
	wxClient  *wechat.Client
	userRepo  repository.UserRepository
	jwtSecret string
}

func NewAuthService(wxClient *wechat.Client, userRepo repository.UserRepository, jwtSecret string) AuthService {
	return &authService{
		wxClient:  wxClient,
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
	}
}

func (s *authService) WxLogin(req *WxLoginReq) (*WxLoginResp, error) {
	logger.Info("auth_service", "WxLogin 开始", map[string]any{
		"nick_name": req.NickName,
		"country":   req.Country,
		"province":  req.Province,
		"city":      req.City,
	})

	// 1. 用 code 换取 openid + session_key
	session, err := s.wxClient.Code2Session(req.Code)
	if err != nil {
		logger.Error("auth_service", "Code2Session 失败", map[string]any{
			"nick_name": req.NickName,
			"error":     err.Error(),
		})
		return nil, err
	}

	// 2. 写入 / 更新用户
	user := &model.User{
		OpenID:     session.OpenID,
		UnionID:    session.UnionID,
		SessionKey: session.SessionKey,
		NickName:   req.NickName,
		AvatarURL:  req.AvatarURL,
		Gender:     req.Gender,
		Country:    req.Country,
		Province:   req.Province,
		City:       req.City,
	}

	if err := s.userRepo.Upsert(user); err != nil {
		logger.Error("auth_service", "用户 Upsert 失败", map[string]any{
			"openid":    session.OpenID,
			"nick_name": req.NickName,
			"error":     err.Error(),
		})
		return nil, err
	}

	logger.Info("auth_service", "用户 Upsert 成功", map[string]any{
		"user_id":   user.ID,
		"openid":    user.OpenID,
		"nick_name": user.NickName,
		"gender":    user.Gender,
		"country":   user.Country,
		"province":  user.Province,
		"city":      user.City,
	})

	// 3. 签发 JWT
	token, err := jwt.Generate(s.jwtSecret, user.ID, user.OpenID)
	if err != nil {
		logger.Error("auth_service", "JWT 签发失败", map[string]any{
			"user_id":   user.ID,
			"openid":    user.OpenID,
			"nick_name": user.NickName,
			"error":     err.Error(),
		})
		return nil, err
	}

	logger.Info("auth_service", "WxLogin 成功，JWT 已签发", map[string]any{
		"user_id":   user.ID,
		"openid":    user.OpenID,
		"nick_name": user.NickName,
	})

	return &WxLoginResp{Token: token, User: user}, nil
}
