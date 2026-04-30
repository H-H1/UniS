package service

import (
	"uniS/internal/model"
	"uniS/internal/repository"
	"uniS/pkg/jwt"
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
	// 1. 用 code 换取 openid + session_key
	session, err := s.wxClient.Code2Session(req.Code)
	if err != nil {
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
		return nil, err
	}

	// 3. 签发 JWT
	token, err := jwt.Generate(s.jwtSecret, user.ID, user.OpenID)
	if err != nil {
		return nil, err
	}

	return &WxLoginResp{Token: token, User: user}, nil
}
