package service

import (
	"time"

	"scau-daily/internal/dto"
	"scau-daily/internal/model"
	"scau-daily/internal/repository"

	"github.com/google/uuid"
)

// TodayService aggregates data for the "Today" overview page.
type TodayService struct{}

// GetOverview returns a combined view of today's courses, pending todos,
// recent notifications, and summary statistics.
func (s *TodayService) GetOverview(userID uuid.UUID) (*dto.TodayOverview, error) {
	now := time.Now()
	date := now.Format("2006-01-02")

	// Weekday: 1=Monday ... 7=Sunday (matches model convention)
	weekday := int((now.Weekday() + 6) % 7 + 1)

	// Current semester info
	semester, err := repository.GetCurrentSemester()
	if err != nil {
		return nil, err
	}

	currentWeek := calcCurrentWeek(semester.StartDate, now)

	// ---- Today's courses ----
	dayOfWeek := int8(weekday)
	courses, err := repository.GetCoursesByUserAndSemester(userID, semester.ID)
	if err != nil {
		return nil, err
	}

	var todayCourses []dto.CourseDTO
	for _, course := range courses {
		for _, sched := range course.Schedules {
			if sched.DayOfWeek == dayOfWeek && isWeekActive(sched.Weeks, currentWeek) {
				todayCourses = append(todayCourses, toCourseDTO(course))
				break
			}
		}
	}

	// ---- Pending todos ----
	pendingTodos, err := repository.GetPendingTodosByUser(userID)
	if err != nil {
		return nil, err
	}

	todoDTOs := make([]dto.TodoDTO, len(pendingTodos))
	for i, t := range pendingTodos {
		todoDTOs[i] = toTodoDTO(t)
	}

	// ---- Recent notifications (last 5) ----
	notifications, err := repository.GetRecentNotifications(5)
	if err != nil {
		return nil, err
	}

	notifDTOs := make([]dto.NotificationDTO, len(notifications))
	for i, n := range notifications {
		notifDTOs[i] = toNotificationDTO(n)
	}

	// ---- Statistics ----
	pendingCount, err := repository.CountPendingTodos(userID)
	if err != nil {
		return nil, err
	}

	daysToFinals := int(semester.EndDate.Sub(now).Hours() / 24)
	if daysToFinals < 0 {
		daysToFinals = 0
	}

	// Weather is a placeholder until a weather API is integrated
	weather := &dto.WeatherInfo{
		Temp:      0,
		Condition: "暂未接入",
		Icon:      "unknown",
	}

	return &dto.TodayOverview{
		Date:          date,
		Weekday:       weekday,
		CurrentWeek:   currentWeek,
		Semester:      semester.Name,
		Weather:       weather,
		Courses:       todayCourses,
		Todos:         todoDTOs,
		Notifications: notifDTOs,
		Stats: dto.TodayStats{
			CoursesToday:  len(todayCourses),
			PendingTodos:  pendingCount,
			DaysToFinals:  daysToFinals,
		},
	}, nil
}

// toNotificationDTO converts a model.Notification to its DTO representation.
func toNotificationDTO(n model.Notification) dto.NotificationDTO {
	return dto.NotificationDTO{
		ID:          n.ID,
		Title:       n.Title,
		Content:     n.Content,
		Source:      n.Source,
		Category:    n.Category,
		PublishedAt: n.PublishedAt,
		AISummary:   n.AISummary,
		AIDeadline:  n.AIDeadline,
		RawURL:      n.RawURL,
	}
}
