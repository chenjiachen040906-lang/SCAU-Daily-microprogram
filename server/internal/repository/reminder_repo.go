package repository

import (
	"time"

	"scau-daily/internal/database"
	"scau-daily/internal/model"

	"github.com/google/uuid"
)

func GetRemindersByUser(userID uuid.UUID) ([]model.Reminder, error) {
	var reminders []model.Reminder
	err := database.DB.Where("user_id = ?", userID).Order("remind_at ASC").Find(&reminders).Error
	return reminders, err
}

func CreateReminder(reminder *model.Reminder) error {
	return database.DB.Create(reminder).Error
}

func UpdateReminder(reminder *model.Reminder) error {
	return database.DB.Save(reminder).Error
}

func DeleteReminder(id uuid.UUID) error {
	return database.DB.Where("id = ?", id).Delete(&model.Reminder{}).Error
}

func FindPendingReminders() ([]model.Reminder, error) {
	var reminders []model.Reminder
	err := database.DB.
		Where("is_notified = ? AND is_enabled = ? AND remind_at <= ?", false, true, time.Now()).
		Order("remind_at ASC").
		Find(&reminders).Error
	return reminders, err
}
