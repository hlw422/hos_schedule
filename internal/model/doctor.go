package model

import "time"

type Doctor struct {
	ID           int64     `gorm:"primaryKey" json:"id"`
	UserID       int64     `gorm:"index" json:"user_id,omitempty"`
	DepartmentID int64     `gorm:"index" json:"department_id"`
	Name         string    `gorm:"size:50" json:"name"`
	Avatar       string    `gorm:"size:255" json:"avatar,omitempty"`
	Title        string    `gorm:"size:50" json:"title,omitempty"`
	Intro        string    `gorm:"type:text" json:"intro,omitempty"`
	Specialty    string    `gorm:"type:text" json:"specialty,omitempty"`
	Status       int8      `gorm:"default:1" json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
