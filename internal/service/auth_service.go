package service

import (
	"fmt"
	"uniS/internal/model"
	"uniS/internal/repository"
	"uniS/pkg/jwt"
	"uniS/pkg/logger"
	"uniS/pkg/wechat"
)

type WxLoginReq struct {
	Code      string `json:"code"` // 本地开发用：wx.login() 获取的 code（云托管时不需要）
	OpenID    string `json:"-"`    // 云托管：由 X-WX-OPENID 请求头注入，handler 中赋值
	UnionID   string `json:"-"`    // 云托管：由 X-WX-UNIONID 请求头注入，handler 中赋值
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

	// 1. 获取 openid
	//    云托管模式：X-WX-OPENID 已由微信平台注入请求头，无需 code 换取
	//    本地开发：用 wx.login() 的 code 调 jscode2session 换取
	openID := req.OpenID
	unionID := req.UnionID
	sessionKey := ""

	if openID == "" {
		if req.Code == "" {
			return nil, fmt.Errorf("缺少 openID（云托管请求头 X-WX-OPENID）或 code（本地开发 wx.login）")
		}
		session, err := s.wxClient.Code2Session(req.Code)
		if err != nil {
			logger.Error("auth_service", "Code2Session 失败", map[string]any{
				"nick_name": req.NickName,
				"error":     err.Error(),
			})
			return nil, err
		}
		openID = session.OpenID
		unionID = session.UnionID
		sessionKey = session.SessionKey
	}

	// 2. 写入 / 更新用户
	user := &model.User{
		OpenID:     openID,
		UnionID:    unionID,
		SessionKey: sessionKey,
		NickName:   req.NickName,
		AvatarURL:  req.AvatarURL,
		Gender:     req.Gender,
		Country:    req.Country,
		Province:   req.Province,
		City:       req.City,
	}

	if err := s.userRepo.Upsert(user); err != nil {
		logger.Error("auth_service", "用户 Upsert 失败", map[string]any{
			"openid":    openID,
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
