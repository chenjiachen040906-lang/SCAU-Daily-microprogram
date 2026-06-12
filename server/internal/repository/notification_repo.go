package repository

import (
	"scau-daily/internal/database"
	"scau-daily/internal/model"

	"github.com/google/uuid"
)

func GetNotifications(page, size int) ([]model.Notification, int64, error) {
	var notifications []model.Notification
	var total int64

	database.DB.Model(&model.Notification{}).Count(&total)

	offset := (page - 1) * size
	err := database.DB.
		Order("published_at DESC").
		Offset(offset).Limit(size).
		Find(&notifications).Error

	return notifications, total, err
}

func GetNotificationByID(id uuid.UUID) (*model.Notification, error) {
	var notification model.Notification
	if err := database.DB.Where("id = ?", id).First(&notification).Error; err != nil {
		return nil, err
	}
	return &notification, nil
}

func CreateNotification(notification *model.Notification) error {
	return database.DB.Create(notification).Error
}

func GetRecentNotifications(limit int) ([]model.Notification, error) {
	var notifications []model.Notification
	err := database.DB.
		Order("published_at DESC").
		Limit(limit).
		Find(&notifications).Error
	return notifications, err
}
