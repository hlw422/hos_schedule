package repository

import (
	"hos_schedule/internal/model"

	"gorm.io/gorm"
)

type HospitalRepo struct {
	db *gorm.DB
}

func NewHospitalRepo(db *gorm.DB) *HospitalRepo {
	return &HospitalRepo{db: db}
}

func (r *HospitalRepo) List() ([]model.Hospital, error) {
	var hospitals []model.Hospital
	err := r.db.Where("status = ?", 1).Find(&hospitals).Error
	return hospitals, err
}

func (r *HospitalRepo) GetByID(id int64) (*model.Hospital, error) {
	var hospital model.Hospital
	err := r.db.First(&hospital, id).Error
	return &hospital, err
}

func (r *HospitalRepo) GetCampuses(hospitalID int64) ([]model.Campus, error) {
	var campuses []model.Campus
	err := r.db.Where("hospital_id = ? AND status = ?", hospitalID, 1).Find(&campuses).Error
	return campuses, err
}

func (r *HospitalRepo) GetCampusByID(id int64) (*model.Campus, error) {
	var campus model.Campus
	err := r.db.First(&campus, id).Error
	return &campus, err
}

func (r *HospitalRepo) GetAllCampuses() ([]model.Campus, error) {
	var campuses []model.Campus
	err := r.db.Where("status = ?", 1).Find(&campuses).Error
	return campuses, err
}

func (r *HospitalRepo) Update(hospital *model.Hospital) error {
	return r.db.Save(hospital).Error
}

func (r *HospitalRepo) CreateCampus(campus *model.Campus) error {
	return r.db.Create(campus).Error
}
