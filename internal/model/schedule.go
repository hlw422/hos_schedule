package model

import "time"

type Schedule struct {
	ID          int64     `gorm:"primaryKey" json:"id"`
	DoctorID    int64     `gorm:"index" json:"doctor_id"`
	CampusID    int64     `gorm:"index" json:"campus_id"`
	Date        string    `gorm:"type:date" json:"date"`
	TimePeriod  string    `gorm:"size:20" json:"time_period"`
	TotalCount  int       `json:"total_count"`
	UsedCount   int       `gorm:"default:0" json:"used_count"`
	RemainCount int       `json:"remain_count"`
	Status      int8      `gorm:"default:1" json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
