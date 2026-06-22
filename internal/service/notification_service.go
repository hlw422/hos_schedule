package service

import (
	"hos_schedule/internal/model"
	"hos_schedule/internal/pkg/wechat"
	"hos_schedule/internal/repository"
)

type NotificationService struct {
	repo   *repository.NotificationRepo
	wechat *wechat.Client
}

func NewNotificationService(repo *repository.NotificationRepo, wechat *wechat.Client) *NotificationService {
	return &NotificationService{repo: repo, wechat: wechat}
}

func (s *NotificationService) SendAppointmentSuccess(userID int64, appointment *model.Appointment) error {
	notification := &model.Notification{
		UserID:  userID,
		Type:    "APPOINTMENT_SUCCESS",
		Status:  "PENDING",
		Content: "预约成功",
	}

	if err := s.repo.Create(notification); err != nil {
		return err
	}

	// TODO: 调用微信订阅消息

	notification.Status = "SENT"
	return s.repo.UpdateStatus(notification.ID, "SENT")
}

func (s *NotificationService) SendReminder(userID int64, appointment *model.Appointment) error {
	notification := &model.Notification{
		UserID:  userID,
		Type:    "REMINDER_1DAY",
		Status:  "PENDING",
		Content: "预约提醒",
	}

	return s.repo.Create(notification)
}
