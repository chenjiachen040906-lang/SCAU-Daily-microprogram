package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ---- User ----

type User struct {
	ID           uuid.UUID     `gorm:"type:uuid;primaryKey" json:"id"`
	StudentID    string        `gorm:"type:varchar(20);uniqueIndex;not null" json:"student_id"`
	Name         string        `gorm:"type:varchar(50)" json:"name"`
	Department   string        `gorm:"type:varchar(100)" json:"department"`
	Major        string        `gorm:"type:varchar(100)" json:"major"`
	Grade        string        `gorm:"type:varchar(4)" json:"grade"`
	OpenID       string        `gorm:"type:varchar(100);column:openid" json:"-"`
	AvatarURL    string        `gorm:"type:text" json:"avatar_url"`
	PasswordHash string        `gorm:"type:text" json:"-"`
	DeviceTokens string        `gorm:"type:jsonb;default:'[]'" json:"-"`
	CreatedAt    time.Time     `json:"created_at"`
	LastLoginAt  *time.Time    `json:"last_login_at"`
	Settings     *UserSettings `gorm:"foreignKey:UserID" json:"settings,omitempty"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}

// ---- Semester ----

type Semester struct {
	ID         uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name       string    `gorm:"type:varchar(30);uniqueIndex;not null" json:"name"`
	StartDate  time.Time `gorm:"type:date" json:"start_date"`
	EndDate    time.Time `gorm:"type:date" json:"end_date"`
	TotalWeeks int       `json:"total_weeks"`
	IsCurrent  bool      `gorm:"default:false" json:"is_current"`
}

// ---- Course ----

type Course struct {
	ID         uuid.UUID        `gorm:"type:uuid;primaryKey" json:"id"`
	UserID     uuid.UUID        `gorm:"type:uuid;not null;index:idx_courses_user_semester" json:"user_id"`
	SemesterID uint             `gorm:"not null;index:idx_courses_user_semester" json:"semester_id"`
	Name       string           `gorm:"type:varchar(100);not null" json:"name"`
	Teacher    string           `gorm:"type:varchar(50)" json:"teacher"`
	Location   string           `gorm:"type:varchar(100)" json:"location"`
	CourseType string           `gorm:"type:varchar(20)" json:"course_type"`
	Credit     float64          `gorm:"type:decimal(3,1)" json:"credit"`
	ExamType   string           `gorm:"type:varchar(20)" json:"exam_type"`
	Schedules  []CourseSchedule `gorm:"foreignKey:CourseID" json:"schedules,omitempty"`
}

func (c *Course) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}

// ---- CourseSchedule ----

type CourseSchedule struct {
	ID           uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	CourseID     uuid.UUID `gorm:"type:uuid;not null;index:idx_course_schedules_course" json:"course_id"`
	DayOfWeek    int8      `gorm:"not null" json:"day_of_week"`
	StartSection int8      `gorm:"not null" json:"start_section"`
	EndSection   int8      `gorm:"not null" json:"end_section"`
	Weeks        string    `gorm:"type:varchar(200)" json:"weeks"`
}

// ---- Todo ----

type Todo struct {
	ID          uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	UserID      uuid.UUID  `gorm:"type:uuid;not null;index:idx_todos_user_active" json:"user_id"`
	Title       string     `gorm:"type:varchar(200);not null" json:"title"`
	Description string     `gorm:"type:text" json:"description"`
	Deadline    *time.Time `gorm:"type:timestamptz" json:"deadline"`
	IsDone      bool       `gorm:"default:false;index:idx_todos_user_active" json:"is_done"`
	Priority    string     `gorm:"type:varchar(10);default:'medium'" json:"priority"`
	Source      string     `gorm:"type:varchar(20);default:'manual'" json:"source"`
	CreatedAt   time.Time  `json:"created_at"`
}

func (t *Todo) BeforeCreate(tx *gorm.DB) error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return nil
}

// ---- Reminder ----

type Reminder struct {
	ID            uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	UserID        uuid.UUID  `gorm:"type:uuid;not null" json:"user_id"`
	Title         string     `gorm:"type:varchar(200);not null" json:"title"`
	RemindAt      time.Time  `gorm:"type:timestamptz;not null;index:idx_reminders_pending" json:"remind_at"`
	Type          string     `gorm:"type:varchar(20)" json:"type"`
	RelatedID     *uuid.UUID `gorm:"type:uuid" json:"related_id"`
	MinutesBefore int        `gorm:"default:15" json:"minutes_before"`
	IsNotified    bool       `gorm:"default:false;index:idx_reminders_pending" json:"is_notified"`
	IsEnabled     bool       `gorm:"default:true;index:idx_reminders_pending" json:"is_enabled"`
}

func (r *Reminder) BeforeCreate(tx *gorm.DB) error {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	return nil
}

// ---- Notification ----

type Notification struct {
	ID          uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	Title       string     `gorm:"type:varchar(200);not null" json:"title"`
	Content     string     `gorm:"type:text" json:"content"`
	Source      string     `gorm:"type:varchar(50)" json:"source"`
	Category    string     `gorm:"type:varchar(20);index:idx_notifications_category" json:"category"`
	PublishedAt time.Time  `gorm:"index:idx_notifications_category" json:"published_at"`
	AISummary   string     `gorm:"type:text" json:"ai_summary"`
	AIDeadline  *time.Time `gorm:"type:timestamptz" json:"ai_deadline"`
	RawURL      string     `gorm:"type:text" json:"raw_url"`
}

func (n *Notification) BeforeCreate(tx *gorm.DB) error {
	if n.ID == uuid.Nil {
		n.ID = uuid.New()
	}
	return nil
}

// ---- ChatSession ----

type ChatSession struct {
	ID        uuid.UUID     `gorm:"type:uuid;primaryKey" json:"id"`
	UserID    uuid.UUID     `gorm:"type:uuid;not null" json:"user_id"`
	Title     string        `gorm:"type:varchar(100)" json:"title"`
	CreatedAt time.Time     `json:"created_at"`
	Messages  []ChatMessage `gorm:"foreignKey:SessionID" json:"messages,omitempty"`
}

func (cs *ChatSession) BeforeCreate(tx *gorm.DB) error {
	if cs.ID == uuid.Nil {
		cs.ID = uuid.New()
	}
	return nil
}

// ---- ChatMessage ----

type ChatMessage struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	SessionID uuid.UUID `gorm:"type:uuid;not null;index:idx_chat_messages_session" json:"session_id"`
	Role      string    `gorm:"type:varchar(20);not null" json:"role"`
	Content   string    `gorm:"type:text" json:"content"`
	ToolCalls string    `gorm:"type:jsonb" json:"tool_calls,omitempty"`
	CreatedAt time.Time `gorm:"index:idx_chat_messages_session" json:"created_at"`
}

func (cm *ChatMessage) BeforeCreate(tx *gorm.DB) error {
	if cm.ID == uuid.Nil {
		cm.ID = uuid.New()
	}
	return nil
}

// ---- UserSettings ----

type UserSettings struct {
	UserID            uuid.UUID `gorm:"type:uuid;primaryKey" json:"user_id"`
	ReminderMinutes   int       `gorm:"default:15" json:"reminder_minutes"`
	DailyBriefEnabled bool      `gorm:"default:true" json:"daily_brief_enabled"`
	DailyBriefTime    string    `gorm:"type:varchar(5);default:'07:30'" json:"daily_brief_time"`
	Theme             string    `gorm:"type:varchar(10);default:'auto'" json:"theme"`
	CurrentSemesterID *uint     `json:"current_semester_id"`
}
