package wechat

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"hos_schedule/internal/config"
)

type Client struct {
	cfg *config.WechatConfig
}

func NewClient(cfg *config.WechatConfig) *Client {
	return &Client{cfg: cfg}
}

type SubscribeMessage struct {
	ToUser     string                 `json:"touser"`
	TemplateID string                 `json:"template_id"`
	Page       string                 `json:"page"`
	Data       map[string]interface{} `json:"data"`
}

func (c *Client) SendSubscribeMessage(msg *SubscribeMessage) error {
	token, err := c.getAccessToken()
	if err != nil {
		return err
	}

	url := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/message/subscribe/send?access_token=%s", token)

	body, _ := json.Marshal(msg)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var result struct {
		ErrCode int    `json:"errcode"`
		ErrMsg  string `json:"errmsg"`
	}
	json.NewDecoder(resp.Body).Decode(&result)

	if result.ErrCode != 0 {
		return fmt.Errorf("wechat send failed: %s", result.ErrMsg)
	}

	return nil
}

func (c *Client) getAccessToken() (string, error) {
	url := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s",
		c.cfg.AppID, c.cfg.Secret)

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result struct {
		AccessToken string `json:"access_token"`
		ErrCode     int    `json:"errcode"`
		ErrMsg      string `json:"errmsg"`
	}
	json.NewDecoder(resp.Body).Decode(&result)

	if result.ErrCode != 0 {
		return "", fmt.Errorf("get access token failed: %s", result.ErrMsg)
	}

	return result.AccessToken, nil
}
