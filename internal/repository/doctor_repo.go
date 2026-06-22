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

func (r *DoctorRepo) GetByUserID(userID int64) (*model.Doctor, error) {
	var doctor model.Doctor
	err := r.db.Where("user_id = ?", userID).First(&doctor).Error
	return &doctor, err
}

func (r *DoctorRepo) ListRecommended(limit int) ([]model.Doctor, error) {
	var doctors []model.Doctor
	err := r.db.Where("status = ?", 1).Limit(limit).Find(&doctors).Error
	return doctors, err
}

func (r *DoctorRepo) Create(doctor *model.Doctor) error {
	return r.db.Create(doctor).Error
}

func (r *DoctorRepo) Update(doctor *model.Doctor) error {
	return r.db.Save(doctor).Error
}

func (r *DoctorRepo) UpdateStatus(id int64, status int8) error {
	return r.db.Model(&model.Doctor{}).Where("id = ?", id).Update("status", status).Error
}
