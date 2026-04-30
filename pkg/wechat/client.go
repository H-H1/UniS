package wechat

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const code2SessionURL = "https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code"

type Code2SessionResp struct {
	OpenID     string `json:"openid"`
	SessionKey string `json:"session_key"`
	UnionID    string `json:"unionid"`
	ErrCode    int    `json:"errcode"`
	ErrMsg     string `json:"errmsg"`
}

type Client struct {
	AppID     string
	AppSecret string
}

func NewClient(appID, appSecret string) *Client {
	return &Client{AppID: appID, AppSecret: appSecret}
}

func (c *Client) Code2Session(code string) (*Code2SessionResp, error) {
	url := fmt.Sprintf(code2SessionURL, c.AppID, c.AppSecret, code)
	resp, err := http.Get(url) //nolint:gosec
	if err != nil {
		return nil, fmt.Errorf("wx request failed: %w", err)
	}
	defer resp.Body.Close()

	var result Code2SessionResp
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("wx response decode failed: %w", err)
	}
	if result.ErrCode != 0 {
		return nil, fmt.Errorf("wx error %d: %s", result.ErrCode, result.ErrMsg)
	}
	return &result, nil
}
