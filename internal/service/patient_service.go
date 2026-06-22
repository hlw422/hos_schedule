package service

import (
	"hos_schedule/internal/model"
	"hos_schedule/internal/repository"
)

type PatientService struct {
	repo *repository.PatientRepo
}

func NewPatientService(repo *repository.PatientRepo) *PatientService {
	return &PatientService{repo: repo}
}

func (s *PatientService) ListByUser(userID int64) ([]model.Patient, error) {
	return s.repo.ListByUser(userID)
}

func (s *PatientService) GetByID(id int64) (*model.Patient, error) {
	return s.repo.GetByID(id)
}

func (s *PatientService) Create(patient *model.Patient) error {
	return s.repo.Create(patient)
}

func (s *PatientService) Update(patient *model.Patient) error {
	return s.repo.Update(patient)
}

func (s *PatientService) Delete(id int64) error {
	return s.repo.Delete(id)
}

func (s *PatientService) SetDefault(userID, patientID int64) error {
	return s.repo.SetDefault(userID, patientID)
}
