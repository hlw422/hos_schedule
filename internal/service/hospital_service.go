package service

import (
	"hos_schedule/internal/model"
	"hos_schedule/internal/repository"
)

type HospitalService struct {
	repo *repository.HospitalRepo
}

func NewHospitalService(repo *repository.HospitalRepo) *HospitalService {
	return &HospitalService{repo: repo}
}

func (s *HospitalService) List() ([]model.Hospital, error) {
	return s.repo.List()
}

func (s *HospitalService) GetByID(id int64) (*model.Hospital, error) {
	return s.repo.GetByID(id)
}

func (s *HospitalService) GetCampuses(hospitalID int64) ([]model.Campus, error) {
	return s.repo.GetCampuses(hospitalID)
}

func (s *HospitalService) GetCampusByID(id int64) (*model.Campus, error) {
	return s.repo.GetCampusByID(id)
}

func (s *HospitalService) GetAllCampuses() ([]model.Campus, error) {
	return s.repo.GetAllCampuses()
}

func (s *HospitalService) Update(hospital *model.Hospital) error {
	return s.repo.Update(hospital)
}

func (s *HospitalService) CreateCampus(campus *model.Campus) error {
	return s.repo.CreateCampus(campus)
}
