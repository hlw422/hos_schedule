package model

import "time"

type Patient struct {
	ID        int64     `gorm:"primaryKey" json:"id"`
	UserID    int64     `gorm:"index" json:"user_id"`
	Name      string    `gorm:"size:50" json:"name"`
	IDCard    string    `gorm:"size:18" json:"id_card,omitempty"`
	Phone     string    `gorm:"size:20" json:"phone,omitempty"`
	Relation  string    `gorm:"size:20" json:"relation,omitempty"`
	IsDefault bool      `gorm:"default:false" json:"is_default"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
