package model

import "time"

type Notification struct {
	ID         int64     `gorm:"primaryKey" json:"id"`
	UserID     int64     `gorm:"index" json:"user_id"`
	Type       string    `gorm:"size:50" json:"type"`
	TemplateID string    `gorm:"size:100" json:"template_id,omitempty"`
	Content    string    `gorm:"type:text" json:"content,omitempty"`
	Status     string    `gorm:"size:20" json:"status"`
	CreatedAt  time.Time `json:"created_at"`
}
