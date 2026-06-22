package service

import (
	"context"
	"fmt"

	"hos_schedule/internal/model"
	redisutil "hos_schedule/internal/pkg/redis"
	"hos_schedule/internal/repository"

	"gorm.io/gorm"
)

type AppointmentService struct {
	db              *gorm.DB
	appointmentRepo *repository.AppointmentRepo
	scheduleRepo    *repository.ScheduleRepo
	slotManager     *redisutil.SlotManager
}

func NewAppointmentService(
	db *gorm.DB,
	appointmentRepo *repository.AppointmentRepo,
	scheduleRepo *repository.ScheduleRepo,
	slotManager *redisutil.SlotManager,
) *AppointmentService {
	return &AppointmentService{
		db:              db,
		appointmentRepo: appointmentRepo,
		scheduleRepo:    scheduleRepo,
		slotManager:     slotManager,
	}
}

type CreateAppointmentRequest struct {
	PatientID  int64  `json:"patient_id" binding:"required"`
	DoctorID   int64  `json:"doctor_id" binding:"required"`
	ScheduleID int64  `json:"schedule_id" binding:"required"`
	CampusID   int64  `json:"campus_id" binding:"required"`
	Date       string `json:"date" binding:"required"`
	TimePeriod string `json:"time_period" binding:"required"`
	PayType    string `json:"pay_type"`
}

func (s *AppointmentService) Create(ctx context.Context, userID int64, req *CreateAppointmentRequest) (*model.Appointment, error) {
	exists, err := s.appointmentRepo.Exists(userID, req.DoctorID, req.Date, req.TimePeriod)
	if err != nil {
		return nil, fmt.Errorf("failed to check duplicate: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("already have appointment for this slot")
	}

	success, err := s.slotManager.DeductSlot(ctx, req.ScheduleID)
	if err != nil {
		if err := s.scheduleRepo.DecrementRemain(req.ScheduleID); err != nil {
			return nil, fmt.Errorf("no remaining slots")
		}
	}
	if !success {
		return nil, fmt.Errorf("no remaining slots")
	}

	appointment := &model.Appointment{
		PatientID:  req.PatientID,
		DoctorID:   req.DoctorID,
		ScheduleID: req.ScheduleID,
		CampusID:   req.CampusID,
		Date:       req.Date,
		TimePeriod: req.TimePeriod,
		PayType:    req.PayType,
		Status:     "PENDING_PAY",
	}

	if err := s.appointmentRepo.Create(appointment); err != nil {
		s.slotManager.ReleaseSlot(ctx, req.ScheduleID)
		return nil, fmt.Errorf("failed to create appointment: %w", err)
	}

	return appointment, nil
}

func (s *AppointmentService) Cancel(ctx context.Context, id int64, reason string) error {
	appointment, err := s.appointmentRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("appointment not found")
	}

	if appointment.Status == "CANCELLED" {
		return fmt.Errorf("appointment already cancelled")
	}

	if err := s.appointmentRepo.UpdateStatus(id, "CANCELLED"); err != nil {
		return err
	}

	if err := s.slotManager.ReleaseSlot(ctx, appointment.ScheduleID); err != nil {
		s.scheduleRepo.IncrementRemain(appointment.ScheduleID)
	}

	return nil
}

func (s *AppointmentService) GetByID(id int64) (*model.Appointment, error) {
	return s.appointmentRepo.GetByID(id)
}

func (s *AppointmentService) ListByUser(userID int64, status string) ([]model.Appointment, error) {
	return s.appointmentRepo.ListByUser(userID, status)
}

func (s *AppointmentService) ListByDoctor(doctorID int64, date string) ([]model.Appointment, error) {
	return s.appointmentRepo.ListByDoctor(doctorID, date)
}
