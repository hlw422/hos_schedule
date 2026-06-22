package repository

import (
	"hos_schedule/internal/model"

	"gorm.io/gorm"
)

type DoctorRepo struct {
	db *gorm.DB
}

func NewDoctorRepo(db *gorm.DB) *DoctorRepo {
	return &DoctorRepo{db: db}
}

func (r *DoctorRepo) ListByDepartment(departmentID int64) ([]model.Doctor, error) {
	var doctors []model.Doctor
	err := r.db.Where("department_id = ? AND status = ?", departmentID, 1).Find(&doctors).Error
	return doctors, err
}

func (r *DoctorRepo) GetByID(id int64) (*model.Doctor, error) {
	var doctor model.Doctor
	err := r.db.First(&doctor, id).Error
	return &doctor, err
}

func (r *DoctorRepo) ListRecommended(limit int) ([]model.Doctor, error) {
	var doctors []model.Doctor
	err := r.db.Where("status = ?", 1).Limit(limit).Find(&doctors).Error
	return doctors, err
}
