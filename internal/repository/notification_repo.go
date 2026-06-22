package repository

import (
	"hos_schedule/internal/model"

	"gorm.io/gorm"
)

type NotificationRepo struct {
	db *gorm.DB
}

func NewNotificationRepo(db *gorm.DB) *NotificationRepo {
	return &NotificationRepo{db: db}
}

func (r *NotificationRepo) Create(notification *model.Notification) error {
	return r.db.Create(notification).Error
}

func (r *NotificationRepo) UpdateStatus(id int64, status string) error {
	return r.db.Model(&model.Notification{}).Where("id = ?", id).
		Update("status", status).Error
}

func (r *NotificationRepo) ListByUser(userID int64) ([]model.Notification, error) {
	var notifications []model.Notification
	err := r.db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&notifications).Error
	return notifications, err
}
