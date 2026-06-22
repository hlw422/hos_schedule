package service

import (
	"hos_schedule/internal/model"
	"hos_schedule/internal/pkg/sms"
	"hos_schedule/internal/pkg/wechat"
	"hos_schedule/internal/repository"
)

type NotificationService struct {
	repo   *repository.NotificationRepo
	wechat *wechat.Client
	sms    *sms.Client
}

func NewNotificationService(repo *repository.NotificationRepo, wechat *wechat.Client) *NotificationService {
	return &NotificationService{
		repo:   repo,
		wechat: wechat,
		sms:    sms.NewClient(),
	}
}

func (s *NotificationService) SendAppointmentSuccess(userID int64, phone, doctorName, date, timePeriod string) error {
	notification := &model.Notification{
		UserID:  userID,
		Type:    "APPOINTMENT_SUCCESS",
		Status:  "PENDING",
		Content: "预约成功",
	}
	if err := s.repo.Create(notification); err != nil {
		return err
	}

	s.sms.SendAppointmentSuccess(phone, doctorName, date, timePeriod)

	return s.repo.UpdateStatus(notification.ID, "SENT")
}

func (s *NotificationService) SendAppointmentReminder(userID int64, phone, doctorName, date, timePeriod string) error {
	notification := &model.Notification{
		UserID:  userID,
		Type:    "REMINDER_1DAY",
		Status:  "PENDING",
		Content: "预约提醒",
	}
	if err := s.repo.Create(notification); err != nil {
		return err
	}

	s.sms.SendAppointmentReminder(phone, doctorName, date, timePeriod)

	return s.repo.UpdateStatus(notification.ID, "SENT")
}

func (s *NotificationService) SendAppointmentCancelled(userID int64, phone, doctorName, date string) error {
	notification := &model.Notification{
		UserID:  userID,
		Type:    "CANCELLED",
		Status:  "PENDING",
		Content: "预约取消",
	}
	if err := s.repo.Create(notification); err != nil {
		return err
	}

	s.sms.SendAppointmentCancelled(phone, doctorName, date)

	return s.repo.UpdateStatus(notification.ID, "SENT")
}
