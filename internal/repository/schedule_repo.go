package repository

import (
	"fmt"

	"hos_schedule/internal/model"

	"gorm.io/gorm"
)

type ScheduleRepo struct {
	db *gorm.DB
}

func NewScheduleRepo(db *gorm.DB) *ScheduleRepo {
	return &ScheduleRepo{db: db}
}

func (r *ScheduleRepo) ListByDoctor(doctorID int64, startDate, endDate string) ([]model.Schedule, error) {
	var schedules []model.Schedule
	err := r.db.Where("doctor_id = ? AND date >= ? AND date <= ? AND status = ?",
		doctorID, startDate, endDate, 1).
		Order("date ASC").
		Find(&schedules).Error
	return schedules, err
}

func (r *ScheduleRepo) ListByDepartment(departmentID int64, date string) ([]model.Schedule, error) {
	var schedules []model.Schedule
	err := r.db.Joins("JOIN doctors ON doctors.id = schedules.doctor_id").
		Where("doctors.department_id = ? AND schedules.date = ? AND schedules.status = ?",
			departmentID, date, 1).
		Find(&schedules).Error
	return schedules, err
}

func (r *ScheduleRepo) GetByID(id int64) (*model.Schedule, error) {
	var schedule model.Schedule
	err := r.db.First(&schedule, id).Error
	return &schedule, err
}

func (r *ScheduleRepo) DecrementRemain(id int64) error {
	result := r.db.Model(&model.Schedule{}).
		Where("id = ? AND remain_count > 0", id).
		Updates(map[string]interface{}{
			"remain_count": gorm.Expr("remain_count - 1"),
			"used_count":   gorm.Expr("used_count + 1"),
		})

	if result.RowsAffected == 0 {
		return fmt.Errorf("no remaining slots")
	}
	return result.Error
}

func (r *ScheduleRepo) Create(schedule *model.Schedule) error {
	return r.db.Create(schedule).Error
}

func (r *ScheduleRepo) IncrementRemain(id int64) error {
	return r.db.Model(&model.Schedule{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"remain_count": gorm.Expr("remain_count + 1"),
			"used_count":   gorm.Expr("used_count - 1"),
		}).Error
}
