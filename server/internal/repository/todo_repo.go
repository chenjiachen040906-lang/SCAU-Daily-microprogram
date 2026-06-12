package repository

import (
	"scau-daily/internal/database"
	"scau-daily/internal/model"

	"github.com/google/uuid"
)

func GetTodosByUser(userID uuid.UUID) ([]model.Todo, error) {
	var todos []model.Todo
	err := database.DB.Where("user_id = ?", userID).Order("created_at DESC").Find(&todos).Error
	return todos, err
}

func GetPendingTodosByUser(userID uuid.UUID) ([]model.Todo, error) {
	var todos []model.Todo
	err := database.DB.Where("user_id = ? AND is_done = ?", userID, false).
		Order("created_at DESC").Find(&todos).Error
	return todos, err
}

func GetDoneTodosByUser(userID uuid.UUID) ([]model.Todo, error) {
	var todos []model.Todo
	err := database.DB.Where("user_id = ? AND is_done = ?", userID, true).
		Order("created_at DESC").Find(&todos).Error
	return todos, err
}

func GetTodoByID(id uuid.UUID) (*model.Todo, error) {
	var todo model.Todo
	if err := database.DB.Where("id = ?", id).First(&todo).Error; err != nil {
		return nil, err
	}
	return &todo, nil
}

func CreateTodo(todo *model.Todo) (*model.Todo, error) {
	if err := database.DB.Create(todo).Error; err != nil {
		return nil, err
	}
	return todo, nil
}

func UpdateTodo(todo *model.Todo) error {
	return database.DB.Save(todo).Error
}

func DeleteTodo(id, userID uuid.UUID) error {
	return database.DB.Where("id = ? AND user_id = ?", id, userID).Delete(&model.Todo{}).Error
}

func CountPendingTodos(userID uuid.UUID) (int, error) {
	var count int64
	err := database.DB.Model(&model.Todo{}).
		Where("user_id = ? AND is_done = ?", userID, false).
		Count(&count).Error
	return int(count), err
}
