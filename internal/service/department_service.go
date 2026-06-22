package service

import (
	"hos_schedule/internal/model"
	"hos_schedule/internal/repository"
)

type DepartmentService struct {
	repo *repository.DepartmentRepo
}

func NewDepartmentService(repo *repository.DepartmentRepo) *DepartmentService {
	return &DepartmentService{repo: repo}
}

func (s *DepartmentService) ListByCampus(campusID int64) ([]model.Department, error) {
	return s.repo.ListByCampus(campusID)
}

func (s *DepartmentService) ListByHospital(hospitalID int64) ([]model.Department, error) {
	return s.repo.ListByHospital(hospitalID)
}

func (s *DepartmentService) GetByID(id int64) (*model.Department, error) {
	return s.repo.GetByID(id)
}
