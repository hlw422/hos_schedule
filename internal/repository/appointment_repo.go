package repository

import (
	"hos_schedule/internal/model"

	"gorm.io/gorm"
)

type AppointmentRepo struct {
	db *gorm.DB
}

func NewAppointmentRepo(db *gorm.DB) *AppointmentRepo {
	return &AppointmentRepo{db: db}
}

func (r *AppointmentRepo) Create(appointment *model.Appointment) error {
	return r.db.Create(appointment).Error
}

func (r *AppointmentRepo) GetByID(id int64) (*model.Appointment, error) {
	var appointment model.Appointment
	err := r.db.First(&appointment, id).Error
	return &appointment, err
}

func (r *AppointmentRepo) ListByUser(userID int64, status string) ([]model.Appointment, error) {
	var appointments []model.Appointment
	query := r.db.Joins("JOIN patients ON patients.id = appointments.patient_id").
		Where("patients.user_id = ?", userID)

	if status != "" {
		query = query.Where("appointments.status = ?", status)
	}

	err := query.Order("appointments.created_at DESC").Find(&appointments).Error
	return appointments, err
}

func (r *AppointmentRepo) ListByDoctor(doctorID int64, date string) ([]model.Appointment, error) {
	var appointments []model.Appointment
	err := r.db.Where("doctor_id = ? AND date = ?", doctorID, date).
		Order("created_at ASC").
		Find(&appointments).Error
	return appointments, err
}

func (r *AppointmentRepo) UpdateStatus(id int64, status string) error {
	return r.db.Model(&model.Appointment{}).Where("id = ?", id).
		Update("status", status).Error
}

func (r *AppointmentRepo) UpdatePayInfo(id int64, status, payID string, payAmount float64) error {
	return r.db.Model(&model.Appointment{}).Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":     status,
			"pay_id":     payID,
			"pay_amount": payAmount,
		}).Error
}

func (r *AppointmentRepo) GetByPayID(payID string) (*model.Appointment, error) {
	var appointment model.Appointment
	err := r.db.Where("pay_id = ?", payID).First(&appointment).Error
	return &appointment, err
}

func (r *AppointmentRepo) Exists(userID, doctorID int64, date, timePeriod string) (bool, error) {
	var count int64
	err := r.db.Model(&model.Appointment{}).
		Joins("JOIN patients ON patients.id = appointments.patient_id").
		Where("patients.user_id = ? AND appointments.doctor_id = ? AND appointments.date = ? AND appointments.time_period = ? AND appointments.status NOT IN ?",
			userID, doctorID, date, timePeriod, []string{"CANCELLED"}).
		Count(&count).Error
	return count > 0, err
}

func (r *AppointmentRepo) CountByDate(date string) (int64, error) {
	var count int64
	err := r.db.Model(&model.Appointment{}).Where("date = ?", date).Count(&count).Error
	return count, err
}

func (r *AppointmentRepo) CountByDateAndStatus(date, status string) (int64, error) {
	var count int64
	err := r.db.Model(&model.Appointment{}).Where("date = ? AND status = ?", date, status).Count(&count).Error
	return count, err
}

func (r *AppointmentRepo) GetMonthlyTrend(year, month int) ([]struct {
	Date  string `json:"date"`
	Count int64  `json:"count"`
}, error) {
	var results []struct {
		Date  string `json:"date"`
		Count int64  `json:"count"`
	}
	err := r.db.Model(&model.Appointment{}).
		Select("date, count(*) as count").
		Where("EXTRACT(YEAR FROM date) = ? AND EXTRACT(MONTH FROM date) = ?", year, month).
		Group("date").
		Order("date ASC").
		Scan(&results).Error
	return results, err
}

func (r *AppointmentRepo) GetTopDepartments(limit int) ([]struct {
	DepartmentID int64  `json:"department_id"`
	Count        int64  `json:"count"`
}, error) {
	var results []struct {
		DepartmentID int64  `json:"department_id"`
		Count        int64  `json:"count"`
	}
	err := r.db.Model(&model.Appointment{}).
		Select("doctor_id, count(*) as count").
		Group("doctor_id").
		Order("count DESC").
		Limit(limit).
		Scan(&results).Error
	return results, err
}

func (r *AppointmentRepo) GetTopDoctors(limit int) ([]struct {
	DoctorID int64  `json:"doctor_id"`
	Count    int64  `json:"count"`
}, error) {
	var results []struct {
		DoctorID int64  `json:"doctor_id"`
		Count    int64  `json:"count"`
	}
	err := r.db.Model(&model.Appointment{}).
		Select("doctor_id, count(*) as count").
		Group("doctor_id").
		Order("count DESC").
		Limit(limit).
		Scan(&results).Error
	return results, err
}
