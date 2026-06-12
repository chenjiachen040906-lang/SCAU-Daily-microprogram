package repository

import (
	"time"

	"scau-daily/internal/database"
	"scau-daily/internal/model"

	"github.com/google/uuid"
)

func FindUserByStudentID(studentID string) (*model.User, error) {
	var user model.User
	if err := database.DB.Where("student_id = ?", studentID).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func FindUserByID(id uuid.UUID) (*model.User, error) {
	var user model.User
	if err := database.DB.Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func FindUserByOpenID(openID string) (*model.User, error) {
	var user model.User
	if err := database.DB.Where("openid = ?", openID).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func CreateUser(user *model.User) (*model.User, error) {
	if err := database.DB.Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func UpdateUser(user *model.User) error {
	return database.DB.Save(user).Error
}

func UpdateUserPassword(id uuid.UUID, hash string) error {
	return database.DB.Model(&model.User{}).Where("id = ?", id).Update("password_hash", hash).Error
}

func UpdateUserLastLogin(id uuid.UUID, t time.Time) error {
	return database.DB.Model(&model.User{}).Where("id = ?", id).Update("last_login_at", t).Error
}
