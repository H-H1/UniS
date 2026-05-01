package wechat

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"uniS/pkg/logger"
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
	url := fmt.Sprintf(code2SessionURL, c.AppID, c.AppSecret, url.QueryEscape(code))

	logger.Info("wechat", "Code2Session 请求微信接口", map[string]any{
		"appid":    c.AppID,
		"code_len": len(code), // 不记录完整 code，避免泄露
	})

	resp, err := http.Get(url) //nolint:gosec
	if err != nil {
		logger.Error("wechat", "HTTP 请求失败", map[string]any{
			"appid": c.AppID,
			"error": err.Error(),
		})
		return nil, fmt.Errorf("wx request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error("wechat", "读取响应体失败", map[string]any{
			"appid":       c.AppID,
			"http_status": resp.StatusCode,
			"error":       err.Error(),
		})
		return nil, fmt.Errorf("wx response read failed: %w", err)
	}

	logger.Info("wechat", "微信接口原始响应", map[string]any{
		"appid":       c.AppID,
		"http_status": resp.StatusCode,
		"body":        string(body),
	})

	var result Code2SessionResp
	if err := json.Unmarshal(body, &result); err != nil {
		logger.Error("wechat", "响应 JSON 解析失败", map[string]any{
			"appid":       c.AppID,
			"http_status": resp.StatusCode,
			"body":        string(body),
			"error":       err.Error(),
		})
		return nil, fmt.Errorf("wx response decode failed (body: %s): %w", string(body), err)
	}

	if result.ErrCode != 0 {
		logger.Error("wechat", "微信返回业务错误", map[string]any{
			"appid":   c.AppID,
			"errcode": result.ErrCode,
			"errmsg":  result.ErrMsg,
			"openid":  result.OpenID, // 通常为空
		})
		return nil, fmt.Errorf("wx error %d: %s", result.ErrCode, result.ErrMsg)
	}

	logger.Info("wechat", "Code2Session 成功", map[string]any{
		"appid":  c.AppID,
		"openid": result.OpenID,
	})
	return &result, nil
}
