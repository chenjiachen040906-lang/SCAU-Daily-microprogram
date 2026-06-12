package handler

import (
	"scau-daily/internal/dto"
	"scau-daily/internal/middleware"
	"scau-daily/internal/service"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// TodoHandler handles todo-related HTTP requests.
type TodoHandler struct {
	svc *service.TodoService
}

// NewTodoHandler creates a new TodoHandler with the given TodoService.
func NewTodoHandler(svc *service.TodoService) *TodoHandler {
	return &TodoHandler{svc: svc}
}

// List returns the user's todos, optionally filtered by status.
// GET /api/todos?status=pending
func (h *TodoHandler) List(c *fiber.Ctx) error {
	userID := middleware.GetUser(c)
	if userID.String() == "00000000-0000-0000-0000-000000000000" {
		return c.Status(fiber.StatusUnauthorized).JSON(dto.ErrorResponse{
			Error: "unauthorized",
		})
	}

	status := c.Query("status")

	todos, err := h.svc.ListTodos(userID, status)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
			Error:   "fetch_failed",
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"todos": todos,
	})
}

// Create creates a new todo for the authenticated user.
// POST /api/todos
func (h *TodoHandler) Create(c *fiber.Ctx) error {
	userID := middleware.GetUser(c)
	if userID.String() == "00000000-0000-0000-0000-000000000000" {
		return c.Status(fiber.StatusUnauthorized).JSON(dto.ErrorResponse{
			Error: "unauthorized",
		})
	}

	var req dto.CreateTodoRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error:   "invalid_request",
			Message: "Failed to parse request body",
		})
	}

	todo, err := h.svc.CreateTodo(userID, &req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
			Error:   "create_failed",
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(todo)
}

// Update updates an existing todo by ID.
// PUT /api/todos/:id
func (h *TodoHandler) Update(c *fiber.Ctx) error {
	userID := middleware.GetUser(c)
	if userID.String() == "00000000-0000-0000-0000-000000000000" {
		return c.Status(fiber.StatusUnauthorized).JSON(dto.ErrorResponse{
			Error: "unauthorized",
		})
	}

	todoID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error:   "invalid_param",
			Message: "Invalid todo ID format",
		})
	}

	var req dto.UpdateTodoRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error:   "invalid_request",
			Message: "Failed to parse request body",
		})
	}

	todo, err := h.svc.UpdateTodo(todoID, userID, &req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
			Error:   "update_failed",
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(todo)
}

// Delete deletes a todo by ID.
// DELETE /api/todos/:id
func (h *TodoHandler) Delete(c *fiber.Ctx) error {
	userID := middleware.GetUser(c)
	if userID.String() == "00000000-0000-0000-0000-000000000000" {
		return c.Status(fiber.StatusUnauthorized).JSON(dto.ErrorResponse{
			Error: "unauthorized",
		})
	}

	todoID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error:   "invalid_param",
			Message: "Invalid todo ID format",
		})
	}

	if err := h.svc.DeleteTodo(todoID, userID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
			Error:   "delete_failed",
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "todo deleted successfully",
	})
}
