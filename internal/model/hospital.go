package model

import "time"

type Hospital struct {
	ID        int64     `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"size:100" json:"name"`
	Address   string    `gorm:"size:255" json:"address,omitempty"`
	Phone     string    `gorm:"size:20" json:"phone,omitempty"`
	Logo      string    `gorm:"size:255" json:"logo,omitempty"`
	Intro     string    `gorm:"type:text" json:"intro,omitempty"`
	Status    int8      `gorm:"default:1" json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Campus struct {
	ID         int64     `gorm:"primaryKey" json:"id"`
	HospitalID int64     `gorm:"index" json:"hospital_id"`
	Name       string    `gorm:"size:100" json:"name"`
	Address    string    `gorm:"size:255" json:"address,omitempty"`
	Phone      string    `gorm:"size:20" json:"phone,omitempty"`
	Latitude   float64   `gorm:"type:decimal(10,7)" json:"latitude,omitempty"`
	Longitude  float64   `gorm:"type:decimal(10,7)" json:"longitude,omitempty"`
	Status     int8      `gorm:"default:1" json:"status"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
