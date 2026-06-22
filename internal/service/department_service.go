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

func (s *DepartmentService) Create(dept *model.Department) error {
	return s.repo.Create(dept)
}

func (s *DepartmentService) Update(dept *model.Department) error {
	return s.repo.Update(dept)
}

func (s *DepartmentService) Delete(id int64) error {
	return s.repo.Delete(id)
}
