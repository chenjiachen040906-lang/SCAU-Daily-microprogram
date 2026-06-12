package handler

import (
	"scau-daily/internal/dto"
	"scau-daily/internal/middleware"
	"scau-daily/internal/service"

	"github.com/gofiber/fiber/v2"
)

// TodayHandler handles the daily overview HTTP requests.
type TodayHandler struct {
	svc *service.TodayService
}

// NewTodayHandler creates a new TodayHandler with the given TodayService.
func NewTodayHandler(svc *service.TodayService) *TodayHandler {
	return &TodayHandler{svc: svc}
}

// Overview returns the aggregated daily overview for the authenticated user.
// GET /api/today/overview
func (h *TodayHandler) Overview(c *fiber.Ctx) error {
	userID := middleware.GetUser(c)
	if userID.String() == "00000000-0000-0000-0000-000000000000" {
		return c.Status(fiber.StatusUnauthorized).JSON(dto.ErrorResponse{
			Error: "unauthorized",
		})
	}

	overview, err := h.svc.GetOverview(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
			Error:   "fetch_failed",
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(overview)
}
