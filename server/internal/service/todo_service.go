package service

import (
	"errors"

	"scau-daily/internal/dto"
	"scau-daily/internal/model"
	"scau-daily/internal/repository"

	"github.com/google/uuid"
)

// TodoService handles to-do item CRUD operations.
type TodoService struct{}

// ListTodos returns the user's to-do items, optionally filtered by status.
// Accepted status values: "pending", "done", "all" (default).
func (s *TodoService) ListTodos(userID uuid.UUID, status string) ([]dto.TodoDTO, error) {
	var todos []model.Todo
	var err error

	switch status {
	case "pending":
		todos, err = repository.GetPendingTodosByUser(userID)
	case "done":
		todos, err = repository.GetDoneTodosByUser(userID)
	default:
		todos, err = repository.GetTodosByUser(userID)
	}

	if err != nil {
		return nil, err
	}

	result := make([]dto.TodoDTO, len(todos))
	for i, t := range todos {
		result[i] = toTodoDTO(t)
	}

	return result, nil
}

// CreateTodo creates a new to-do item for the user.
func (s *TodoService) CreateTodo(userID uuid.UUID, req *dto.CreateTodoRequest) (*dto.TodoDTO, error) {
	if req.Title == "" {
		return nil, errors.New("标题不能为空")
	}

	priority := req.Priority
	if priority == "" {
		priority = "medium"
	}

	todo := &model.Todo{
		UserID:      userID,
		Title:       req.Title,
		Description: req.Description,
		Deadline:    req.Deadline,
		IsDone:      false,
		Priority:    priority,
		Source:      "manual",
	}

	created, err := repository.CreateTodo(todo)
	if err != nil {
		return nil, err
	}

	d := toTodoDTO(*created)
	return &d, nil
}

// UpdateTodo partially updates a to-do item. Only non-nil fields in the request
// are applied. Returns an error if the to-do does not belong to the user.
func (s *TodoService) UpdateTodo(id uuid.UUID, userID uuid.UUID, req *dto.UpdateTodoRequest) (*dto.TodoDTO, error) {
	todo, err := repository.GetTodoByID(id)
	if err != nil {
		return nil, errors.New("待办不存在")
	}
	if todo.UserID != userID {
		return nil, errors.New("无权操作")
	}

	// Apply partial updates
	if req.Title != nil {
		todo.Title = *req.Title
	}
	if req.Description != nil {
		todo.Description = *req.Description
	}
	if req.Deadline != nil {
		todo.Deadline = req.Deadline
	}
	if req.IsDone != nil {
		todo.IsDone = *req.IsDone
	}
	if req.Priority != nil {
		todo.Priority = *req.Priority
	}

	if err := repository.UpdateTodo(todo); err != nil {
		return nil, err
	}

	d := toTodoDTO(*todo)
	return &d, nil
}

// DeleteTodo removes a to-do item. Returns an error if the to-do does not
// belong to the user.
func (s *TodoService) DeleteTodo(id uuid.UUID, userID uuid.UUID) error {
	return repository.DeleteTodo(id, userID)
}

// toTodoDTO converts a model.Todo to its DTO representation.
func toTodoDTO(todo model.Todo) dto.TodoDTO {
	return dto.TodoDTO{
		ID:          todo.ID,
		Title:       todo.Title,
		Description: todo.Description,
		Deadline:    todo.Deadline,
		IsDone:      todo.IsDone,
		Priority:    todo.Priority,
		Source:      todo.Source,
		CreatedAt:   todo.CreatedAt,
	}
}
