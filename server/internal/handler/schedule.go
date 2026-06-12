package handler

import (
	"strconv"
	"time"

	"scau-daily/internal/dto"
	"scau-daily/internal/middleware"
	"scau-daily/internal/service"

	"github.com/gofiber/fiber/v2"
)

// ScheduleHandler handles schedule-related HTTP requests.
type ScheduleHandler struct {
	svc *service.ScheduleService
}

// NewScheduleHandler creates a new ScheduleHandler with the given ScheduleService.
func NewScheduleHandler(svc *service.ScheduleService) *ScheduleHandler {
	return &ScheduleHandler{svc: svc}
}

// Sync triggers a schedule sync from the academic system for the authenticated user.
// POST /api/schedule/sync
func (h *ScheduleHandler) Sync(c *fiber.Ctx) error {
	userID := middleware.GetUser(c)
	if userID.String() == "00000000-0000-0000-0000-000000000000" {
		return c.Status(fiber.StatusUnauthorized).JSON(dto.ErrorResponse{
			Error: "unauthorized",
		})
	}

	var req dto.SyncScheduleRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error:   "invalid_request",
			Message: "Failed to parse request body",
		})
	}

	resp, err := h.svc.SyncSchedule(userID, req.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
			Error:   "sync_failed",
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}

// Today returns the user's courses for a given date.
// GET /api/schedule/today?date=2006-01-02
func (h *ScheduleHandler) Today(c *fiber.Ctx) error {
	userID := middleware.GetUser(c)
	if userID.String() == "00000000-0000-0000-0000-000000000000" {
		return c.Status(fiber.StatusUnauthorized).JSON(dto.ErrorResponse{
			Error: "unauthorized",
		})
	}

	date := c.Query("date", time.Now().Format("2006-01-02"))

	courses, err := h.svc.GetTodayCourses(userID, date)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
			Error:   "fetch_failed",
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"date":    date,
		"courses": courses,
	})
}

// Week returns the user's courses for a given week number.
// GET /api/schedule/week?week=1
func (h *ScheduleHandler) Week(c *fiber.Ctx) error {
	userID := middleware.GetUser(c)
	if userID.String() == "00000000-0000-0000-0000-000000000000" {
		return c.Status(fiber.StatusUnauthorized).JSON(dto.ErrorResponse{
			Error: "unauthorized",
		})
	}

	weekStr := c.Query("week")
	if weekStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error:   "missing_param",
			Message: "Query parameter 'week' is required",
		})
	}

	week, err := strconv.Atoi(weekStr)
	if err != nil || week < 1 {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error:   "invalid_param",
			Message: "Query parameter 'week' must be a positive integer",
		})
	}

	courses, err := h.svc.GetWeekCourses(userID, week)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
			Error:   "fetch_failed",
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.WeekScheduleResponse{
		Week:    week,
		Courses: courses,
	})
}

// Courses returns all courses for the authenticated user.
// GET /api/schedule/courses
func (h *ScheduleHandler) Courses(c *fiber.Ctx) error {
	userID := middleware.GetUser(c)
	if userID.String() == "00000000-0000-0000-0000-000000000000" {
		return c.Status(fiber.StatusUnauthorized).JSON(dto.ErrorResponse{
			Error: "unauthorized",
		})
	}

	courses, err := h.svc.GetAllCourses(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
			Error:   "fetch_failed",
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"courses": courses,
	})
}

// FreeSlots returns the user's free time slots for a given date.
// GET /api/schedule/free-slots?date=2006-01-02
func (h *ScheduleHandler) FreeSlots(c *fiber.Ctx) error {
	userID := middleware.GetUser(c)
	if userID.String() == "00000000-0000-0000-0000-000000000000" {
		return c.Status(fiber.StatusUnauthorized).JSON(dto.ErrorResponse{
			Error: "unauthorized",
		})
	}

	date := c.Query("date", time.Now().Format("2006-01-02"))

	resp, err := h.svc.GetFreeSlots(userID, date)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
			Error:   "fetch_failed",
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}
