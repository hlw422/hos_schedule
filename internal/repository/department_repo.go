package repository

import (
	"hos_schedule/internal/model"

	"gorm.io/gorm"
)

type DepartmentRepo struct {
	db *gorm.DB
}

func NewDepartmentRepo(db *gorm.DB) *DepartmentRepo {
	return &DepartmentRepo{db: db}
}

func (r *DepartmentRepo) ListByCampus(campusID int64) ([]model.Department, error) {
	var departments []model.Department
	err := r.db.Where("campus_id = ? AND status = ?", campusID, 1).
		Order("sort_order ASC").
		Find(&departments).Error
	return departments, err
}

func (r *DepartmentRepo) ListByHospital(hospitalID int64) ([]model.Department, error) {
	var departments []model.Department
	err := r.db.Where("hospital_id = ? AND status = ?", hospitalID, 1).
		Order("sort_order ASC").
		Find(&departments).Error
	return departments, err
}

func (r *DepartmentRepo) GetByID(id int64) (*model.Department, error) {
	var department model.Department
	err := r.db.First(&department, id).Error
	return &department, err
}

func (r *DepartmentRepo) Create(dept *model.Department) error {
	return r.db.Create(dept).Error
}

func (r *DepartmentRepo) Update(dept *model.Department) error {
	return r.db.Save(dept).Error
}

func (r *DepartmentRepo) Delete(id int64) error {
	return r.db.Delete(&model.Department{}, id).Error
}
