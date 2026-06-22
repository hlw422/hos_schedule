package service

import (
	"hos_schedule/internal/model"
	"hos_schedule/internal/repository"
)

type ScheduleService struct {
	repo *repository.ScheduleRepo
}

func NewScheduleService(repo *repository.ScheduleRepo) *ScheduleService {
	return &ScheduleService{repo: repo}
}

func (s *ScheduleService) ListByDoctor(doctorID int64, startDate, endDate string) ([]model.Schedule, error) {
	return s.repo.ListByDoctor(doctorID, startDate, endDate)
}

func (s *ScheduleService) ListByDepartment(departmentID int64, date string) ([]model.Schedule, error) {
	return s.repo.ListByDepartment(departmentID, date)
}

func (s *ScheduleService) Create(schedule *model.Schedule) error {
	return s.repo.Create(schedule)
}

func (s *ScheduleService) GetByID(id int64) (*model.Schedule, error) {
	return s.repo.GetByID(id)
}
