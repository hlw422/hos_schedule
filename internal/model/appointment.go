package model

import "time"

type Appointment struct {
	ID           int64     `gorm:"primaryKey" json:"id"`
	PatientID    int64     `gorm:"index" json:"patient_id"`
	DoctorID     int64     `gorm:"index" json:"doctor_id"`
	ScheduleID   int64     `gorm:"index" json:"schedule_id"`
	CampusID     int64     `json:"campus_id"`
	Date         string    `gorm:"type:date" json:"date"`
	TimePeriod   string    `gorm:"size:20" json:"time_period"`
	Status       string    `gorm:"size:20;default:PENDING_PAY" json:"status"`
	PayType      string    `gorm:"size:20" json:"pay_type,omitempty"`
	PayAmount    float64   `gorm:"type:decimal(10,2)" json:"pay_amount,omitempty"`
	PayID        string    `gorm:"size:100" json:"pay_id,omitempty"`
	CancelReason string    `gorm:"size:255" json:"cancel_reason,omitempty"`
	VisitNo      string    `gorm:"size:50" json:"visit_no,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
