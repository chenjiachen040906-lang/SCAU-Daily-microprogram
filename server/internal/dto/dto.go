package dto

import (
	"time"

	"github.com/google/uuid"
)

// ===== Auth =====

type LoginRequest struct {
	StudentID string `json:"student_id"`
	Password  string `json:"password"`
}

type WxLoginRequest struct {
	Code string `json:"code"`
}

type BindStudentRequest struct {
	StudentID string `json:"student_id"`
	Password  string `json:"password"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type AuthResponse struct {
	AccessToken  string  `json:"access_token"`
	RefreshToken string  `json:"refresh_token"`
	User         UserDTO `json:"user"`
}

type UserDTO struct {
	ID         uuid.UUID `json:"id"`
	StudentID  string    `json:"student_id"`
	Name       string    `json:"name"`
	Department string    `json:"department"`
	Major      string    `json:"major"`
	Grade      string    `json:"grade"`
	AvatarURL  string    `json:"avatar_url"`
}

// ===== Schedule =====

type CourseDTO struct {
	ID         uuid.UUID           `json:"id"`
	Name       string              `json:"name"`
	Teacher    string              `json:"teacher"`
	Location   string              `json:"location"`
	CourseType string              `json:"course_type"`
	Credit     float64             `json:"credit"`
	ExamType   string              `json:"exam_type"`
	Schedules  []CourseScheduleDTO `json:"schedules"`
}

type CourseScheduleDTO struct {
	DayOfWeek    int8   `json:"day_of_week"`
	StartSection int8   `json:"start_section"`
	EndSection   int8   `json:"end_section"`
	Weeks        string `json:"weeks"`
}

type SyncScheduleRequest struct {
	Password string `json:"password"`
}

type SyncScheduleResponse struct {
	Message string      `json:"message"`
	Courses []CourseDTO `json:"courses"`
	Count   int         `json:"count"`
}

type WeekScheduleResponse struct {
	Week    int         `json:"week"`
	Courses []CourseDTO `json:"courses"`
}

type FreeSlotResponse struct {
	Date      string   `json:"date"`
	FreeSlots []string `json:"free_slots"`
}

// ===== Todo =====

type CreateTodoRequest struct {
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Deadline    *time.Time `json:"deadline"`
	Priority    string     `json:"priority"`
}

type UpdateTodoRequest struct {
	Title       *string    `json:"title"`
	Description *string    `json:"description"`
	Deadline    *time.Time `json:"deadline"`
	IsDone      *bool      `json:"is_done"`
	Priority    *string    `json:"priority"`
}

type TodoDTO struct {
	ID          uuid.UUID  `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Deadline    *time.Time `json:"deadline"`
	IsDone      bool       `json:"is_done"`
	Priority    string     `json:"priority"`
	Source      string     `json:"source"`
	CreatedAt   time.Time  `json:"created_at"`
}

// ===== Reminder =====

type CreateReminderRequest struct {
	Title         string     `json:"title"`
	RemindAt      time.Time  `json:"remind_at"`
	Type          string     `json:"type"`
	RelatedID     *uuid.UUID `json:"related_id"`
	MinutesBefore int        `json:"minutes_before"`
}

type UpdateReminderRequest struct {
	Title     *string `json:"title"`
	IsEnabled *bool   `json:"is_enabled"`
}

type ReminderDTO struct {
	ID            uuid.UUID  `json:"id"`
	Title         string     `json:"title"`
	RemindAt      time.Time  `json:"remind_at"`
	Type          string     `json:"type"`
	RelatedID     *uuid.UUID `json:"related_id"`
	MinutesBefore int        `json:"minutes_before"`
	IsNotified    bool       `json:"is_notified"`
	IsEnabled     bool       `json:"is_enabled"`
}

type DeviceTokenRequest struct {
	Token    string `json:"token"`
	Platform string `json:"platform"`
}

// ===== Today =====

type TodayOverview struct {
	Date          string            `json:"date"`
	Weekday       int               `json:"weekday"`
	CurrentWeek   int               `json:"current_week"`
	Semester      string            `json:"semester"`
	Weather       *WeatherInfo      `json:"weather"`
	Courses       []CourseDTO       `json:"courses"`
	Todos         []TodoDTO         `json:"todos"`
	Notifications []NotificationDTO `json:"notifications"`
	Stats         TodayStats        `json:"stats"`
}

type WeatherInfo struct {
	Temp      int    `json:"temp"`
	Condition string `json:"condition"`
	Icon      string `json:"icon"`
}

type TodayStats struct {
	CoursesToday int `json:"courses_today"`
	PendingTodos int `json:"pending_todos"`
	DaysToFinals int `json:"days_to_finals"`
}

// ===== Notification =====

type NotificationDTO struct {
	ID          uuid.UUID  `json:"id"`
	Title       string     `json:"title"`
	Content     string     `json:"content"`
	Source      string     `json:"source"`
	Category    string     `json:"category"`
	PublishedAt time.Time  `json:"published_at"`
	AISummary   string     `json:"ai_summary"`
	AIDeadline  *time.Time `json:"ai_deadline"`
	RawURL      string     `json:"raw_url"`
}

// ===== Chat =====

type CreateChatSessionRequest struct {
	Title string `json:"title"`
}

type ChatSessionDTO struct {
	ID        uuid.UUID `json:"id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
}

type SendMessageRequest struct {
	Content string `json:"content"`
}

// ===== Settings =====

type SettingsDTO struct {
	ReminderMinutes   int    `json:"reminder_minutes"`
	DailyBriefEnabled bool   `json:"daily_brief_enabled"`
	DailyBriefTime    string `json:"daily_brief_time"`
	Theme             string `json:"theme"`
	CurrentSemesterID *uint  `json:"current_semester_id"`
}

type UpdateSettingsRequest struct {
	ReminderMinutes   *int    `json:"reminder_minutes"`
	DailyBriefEnabled *bool   `json:"daily_brief_enabled"`
	DailyBriefTime    *string `json:"daily_brief_time"`
	Theme             *string `json:"theme"`
	CurrentSemesterID *uint   `json:"current_semester_id"`
}

// ===== Common =====

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Page       int         `json:"page"`
	Size       int         `json:"size"`
	TotalCount int64       `json:"total_count"`
}
