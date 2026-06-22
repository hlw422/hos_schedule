package model

import "time"

type User struct {
	ID        int64     `gorm:"primaryKey" json:"id"`
	OpenID    string    `gorm:"uniqueIndex;size:100" json:"openid"`
	UnionID   string    `gorm:"size:100" json:"unionid,omitempty"`
	Phone     string    `gorm:"size:20" json:"phone,omitempty"`
	Nickname  string    `gorm:"size:50" json:"nickname,omitempty"`
	Avatar    string    `gorm:"size:255" json:"avatar,omitempty"`
	Role      string    `gorm:"size:20;default:PATIENT" json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
