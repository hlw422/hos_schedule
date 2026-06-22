package service

import (
	"hos_schedule/internal/model"
	"hos_schedule/internal/repository"
)

type DoctorService struct {
	repo *repository.DoctorRepo
}

func NewDoctorService(repo *repository.DoctorRepo) *DoctorService {
	return &DoctorService{repo: repo}
}

func (s *DoctorService) ListByDepartment(departmentID int64) ([]model.Doctor, error) {
	return s.repo.ListByDepartment(departmentID)
}

func (s *DoctorService) GetByID(id int64) (*model.Doctor, error) {
	return s.repo.GetByID(id)
}

func (s *DoctorService) GetByUserID(userID int64) (*model.Doctor, error) {
	return s.repo.GetByUserID(userID)
}

func (s *DoctorService) ListRecommended(limit int) ([]model.Doctor, error) {
	return s.repo.ListRecommended(limit)
}

func (s *DoctorService) Create(doctor *model.Doctor) error {
	return s.repo.Create(doctor)
}

func (s *DoctorService) Update(doctor *model.Doctor) error {
	return s.repo.Update(doctor)
}

func (s *DoctorService) UpdateStatus(id int64, status int8) error {
	return s.repo.UpdateStatus(id, status)
}
