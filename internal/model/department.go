package model

import "time"

type Department struct {
	ID         int64     `gorm:"primaryKey" json:"id"`
	HospitalID int64     `gorm:"index" json:"hospital_id"`
	CampusID   int64     `gorm:"index" json:"campus_id,omitempty"`
	Name       string    `gorm:"size:100" json:"name"`
	Intro      string    `gorm:"type:text" json:"intro,omitempty"`
	SortOrder  int       `gorm:"default:0" json:"sort_order"`
	Status     int8      `gorm:"default:1" json:"status"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
