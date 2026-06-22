package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"hos_schedule/internal/config"
	"hos_schedule/internal/model"
	"hos_schedule/internal/pkg/jwt"

	"gorm.io/gorm"
)

type AuthService struct {
	db  *gorm.DB
	cfg *config.Config
}

func NewAuthService(db *gorm.DB, cfg *config.Config) *AuthService {
	return &AuthService{db: db, cfg: cfg}
}

type WechatLoginResp struct {
	OpenID     string `json:"openid"`
	UnionID    string `json:"unionid"`
	SessionKey string `json:"session_key"`
	ErrCode    int    `json:"errcode"`
	ErrMsg     string `json:"errmsg"`
}

func (s *AuthService) WechatLogin(code string) (string, *model.User, error) {
	url := fmt.Sprintf(
		"https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code",
		s.cfg.Wechat.AppID, s.cfg.Wechat.Secret, code,
	)

	resp, err := http.Get(url)
	if err != nil {
		return "", nil, fmt.Errorf("failed to call wechat api: %w", err)
	}
	defer resp.Body.Close()

	var wechatResp WechatLoginResp
	if err := json.NewDecoder(resp.Body).Decode(&wechatResp); err != nil {
		return "", nil, fmt.Errorf("failed to decode wechat response: %w", err)
	}

	if wechatResp.ErrCode != 0 {
		return "", nil, fmt.Errorf("wechat login failed: %s", wechatResp.ErrMsg)
	}

	var user model.User
	result := s.db.Where("openid = ?", wechatResp.OpenID).First(&user)

	if result.Error == gorm.ErrRecordNotFound {
		user = model.User{
			OpenID:   wechatResp.OpenID,
			UnionID:  wechatResp.UnionID,
			Role:     "PATIENT",
			Nickname: "微信用户",
		}
		if err := s.db.Create(&user).Error; err != nil {
			return "", nil, fmt.Errorf("failed to create user: %w", err)
		}
	} else if result.Error != nil {
		return "", nil, fmt.Errorf("failed to query user: %w", result.Error)
	}

	expire, _ := time.ParseDuration(s.cfg.JWT.Expire)
	token, err := jwt.GenerateToken(s.cfg.JWT.Secret, user.ID, user.Role, expire)
	if err != nil {
		return "", nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return token, &user, nil
}
