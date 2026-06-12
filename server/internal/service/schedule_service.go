package service

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"scau-daily/internal/config"
	"scau-daily/internal/dto"
	"scau-daily/internal/jwxt"
	"scau-daily/internal/model"
	"scau-daily/internal/repository"

	"github.com/google/uuid"
)

// sectionTimes maps class section numbers to their start/end times.
// Standard SCAU class schedule.
var sectionTimes = map[int8][2]string{
	1:  {"08:00", "08:45"},
	2:  {"08:55", "09:40"},
	3:  {"10:00", "10:45"},
	4:  {"10:55", "11:40"},
	5:  {"14:30", "15:15"},
	6:  {"15:25", "16:10"},
	7:  {"16:30", "17:15"},
	8:  {"17:25", "18:10"},
	9:  {"19:30", "20:15"},
	10: {"20:25", "21:10"},
	11: {"21:20", "22:05"},
}

// ScheduleService handles course schedule operations.
type ScheduleService struct{}

// SyncSchedule fetches the course schedule from JWXT, replaces existing courses
// for the current semester, and returns the synced courses.
// The password parameter is the user's JWXT password (plain text), required for
// authenticating against the educational system. The handler layer is responsible
// for obtaining it from the request body or a secure session store.
func (s *ScheduleService) SyncSchedule(userID uuid.UUID, password string) (*dto.SyncScheduleResponse, error) {
	user, err := repository.FindUserByID(userID)
	if err != nil {
		return nil, errors.New("用户不存在")
	}
	if user.StudentID == "" {
		return nil, errors.New("未绑定学号，请先绑定")
	}

	semester, err := repository.GetCurrentSemester()
	if err != nil {
		return nil, errors.New("获取当前学期失败: " + err.Error())
	}

	// Login to JWXT
	client := jwxt.NewClient(config.AppConfig.JWXTBaseURL)
	result, err := client.Login(user.StudentID, password)
	if err != nil {
		return nil, errors.New("教务系统连接失败: " + err.Error())
	}
	if !result.Success {
		return nil, errors.New("教务系统登录失败: " + result.Message)
	}

	// Fetch schedule from JWXT
	items, err := client.GetSchedule(user.StudentID, semester.Name)
	if err != nil {
		return nil, errors.New("获取课表失败: " + err.Error())
	}

	// Delete existing courses for the current semester
	if err := repository.DeleteCoursesByUserAndSemester(userID, semester.ID); err != nil {
		return nil, errors.New("删除旧课表失败: " + err.Error())
	}

	// Group schedule items by course name (each jwxt item is one time slot)
	type courseGroup struct {
		name       string
		teacher    string
		location   string
		courseType string
		credit     string
		examType   string
		schedules  []model.CourseSchedule
	}

	groups := make(map[string]*courseGroup)
	var order []string

	for _, item := range items {
		g, exists := groups[item.Name]
		if !exists {
			g = &courseGroup{
				name:       item.Name,
				teacher:    item.Teacher,
				location:   item.Location,
				courseType: item.CourseType,
				credit:     item.Credit,
				examType:   item.ExamType,
			}
			groups[item.Name] = g
			order = append(order, item.Name)
		}
		g.schedules = append(g.schedules, model.CourseSchedule{
			DayOfWeek:    int8(item.DayOfWeek),
			StartSection: int8(item.StartSec),
			EndSection:   int8(item.EndSec),
			Weeks:        item.Weeks,
		})
	}

	// Create course records with their schedules
	var courses []model.Course
	for _, name := range order {
		g := groups[name]
		credit, _ := strconv.ParseFloat(g.credit, 64)
		courses = append(courses, model.Course{
			UserID:     userID,
			SemesterID: semester.ID,
			Name:       g.name,
			Teacher:    g.teacher,
			Location:   g.location,
			CourseType: g.courseType,
			Credit:     credit,
			ExamType:   g.examType,
			Schedules:  g.schedules,
		})
	}

	if err := repository.BatchCreateCourses(courses); err != nil {
		return nil, errors.New("保存课表失败: " + err.Error())
	}

	// Convert to DTOs for response
	courseDTOs := make([]dto.CourseDTO, len(courses))
	for i, c := range courses {
		courseDTOs[i] = toCourseDTO(c)
	}

	return &dto.SyncScheduleResponse{
		Message: fmt.Sprintf("同步成功，共 %d 门课程", len(courses)),
		Courses: courseDTOs,
		Count:   len(courses),
	}, nil
}

// GetTodayCourses returns courses scheduled for the given date.
// It filters by day-of-week and checks whether the current week is in each
// schedule's active week list.
func (s *ScheduleService) GetTodayCourses(userID uuid.UUID, date string) ([]dto.CourseDTO, error) {
	t, err := time.Parse("2006-01-02", date)
	if err != nil {
		return nil, errors.New("日期格式错误，请使用 YYYY-MM-DD")
	}

	semester, err := repository.GetCurrentSemester()
	if err != nil {
		return nil, errors.New("获取当前学期失败: " + err.Error())
	}

	courses, err := repository.GetCoursesByUserAndSemester(userID, semester.ID)
	if err != nil {
		return nil, err
	}

	// Determine day of week (1=Monday ... 7=Sunday), matching model convention
	dayOfWeek := int8((t.Weekday() + 6) % 7 + 1)

	// Calculate current week number from semester start date
	currentWeek := calcCurrentWeek(semester.StartDate, t)

	var result []dto.CourseDTO
	for _, course := range courses {
		for _, sched := range course.Schedules {
			if sched.DayOfWeek == dayOfWeek && isWeekActive(sched.Weeks, currentWeek) {
				result = append(result, toCourseDTO(course))
				break
			}
		}
	}

	return result, nil
}

// GetWeekCourses returns courses that have at least one schedule active in the given week.
func (s *ScheduleService) GetWeekCourses(userID uuid.UUID, week int) ([]dto.CourseDTO, error) {
	semester, err := repository.GetCurrentSemester()
	if err != nil {
		return nil, errors.New("获取当前学期失败: " + err.Error())
	}

	courses, err := repository.GetCoursesByUserAndSemester(userID, semester.ID)
	if err != nil {
		return nil, err
	}

	var result []dto.CourseDTO
	for _, course := range courses {
		for _, sched := range course.Schedules {
			if isWeekActive(sched.Weeks, week) {
				result = append(result, toCourseDTO(course))
				break
			}
		}
	}

	return result, nil
}

// GetAllCourses returns every course for the current semester.
func (s *ScheduleService) GetAllCourses(userID uuid.UUID) ([]dto.CourseDTO, error) {
	semester, err := repository.GetCurrentSemester()
	if err != nil {
		return nil, errors.New("获取当前学期失败: " + err.Error())
	}

	courses, err := repository.GetCoursesByUserAndSemester(userID, semester.ID)
	if err != nil {
		return nil, err
	}

	result := make([]dto.CourseDTO, len(courses))
	for i, c := range courses {
		result[i] = toCourseDTO(c)
	}

	return result, nil
}

// GetFreeSlots calculates free time slots for the given date by comparing
// occupied class sections against the full daily schedule (sections 1-11).
func (s *ScheduleService) GetFreeSlots(userID uuid.UUID, date string) (*dto.FreeSlotResponse, error) {
	t, err := time.Parse("2006-01-02", date)
	if err != nil {
		return nil, errors.New("日期格式错误，请使用 YYYY-MM-DD")
	}

	semester, err := repository.GetCurrentSemester()
	if err != nil {
		return nil, errors.New("获取当前学期失败: " + err.Error())
	}

	courses, err := repository.GetCoursesByUserAndSemester(userID, semester.ID)
	if err != nil {
		return nil, err
	}

	dayOfWeek := int8((t.Weekday() + 6) % 7 + 1)
	currentWeek := calcCurrentWeek(semester.StartDate, t)

	// Track which sections are occupied on this day/week
	occupied := make(map[int8]bool)
	for _, course := range courses {
		for _, sched := range course.Schedules {
			if sched.DayOfWeek == dayOfWeek && isWeekActive(sched.Weeks, currentWeek) {
				for sec := sched.StartSection; sec <= sched.EndSection; sec++ {
					occupied[sec] = true
				}
			}
		}
	}

	// Build contiguous free slot ranges
	var freeSlots []string
	var rangeStart int8 = -1

	for sec := int8(1); sec <= 11; sec++ {
		if !occupied[sec] {
			if rangeStart == -1 {
				rangeStart = sec
			}
		} else {
			if rangeStart != -1 {
				startTime := sectionTimes[rangeStart][0]
				endTime := sectionTimes[sec-1][1]
				freeSlots = append(freeSlots, fmt.Sprintf("%s-%s", startTime, endTime))
				rangeStart = -1
			}
		}
	}

	// Handle trailing free sections
	if rangeStart != -1 {
		startTime := sectionTimes[rangeStart][0]
		endTime := sectionTimes[11][1]
		freeSlots = append(freeSlots, fmt.Sprintf("%s-%s", startTime, endTime))
	}

	return &dto.FreeSlotResponse{
		Date:      date,
		FreeSlots: freeSlots,
	}, nil
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

// toCourseDTO converts a model.Course to its DTO representation.
func toCourseDTO(course model.Course) dto.CourseDTO {
	schedules := make([]dto.CourseScheduleDTO, len(course.Schedules))
	for i, s := range course.Schedules {
		schedules[i] = dto.CourseScheduleDTO{
			DayOfWeek:    s.DayOfWeek,
			StartSection: s.StartSection,
			EndSection:   s.EndSection,
			Weeks:        s.Weeks,
		}
	}

	return dto.CourseDTO{
		ID:         course.ID,
		Name:       course.Name,
		Teacher:    course.Teacher,
		Location:   course.Location,
		CourseType: course.CourseType,
		Credit:     course.Credit,
		ExamType:   course.ExamType,
		Schedules:  schedules,
	}
}

// isWeekActive checks whether a given week number appears in the comma-separated
// weeks string (e.g., "1,2,3,4,5,6,7,8").
func isWeekActive(weeks string, week int) bool {
	if weeks == "" {
		return true // empty weeks means always active
	}
	target := strconv.Itoa(week)
	for _, w := range strings.Split(weeks, ",") {
		if strings.TrimSpace(w) == target {
			return true
		}
	}
	return false
}

// calcCurrentWeek computes the week number given a semester start date and a
// reference date. Week 1 starts on the semester start date.
func calcCurrentWeek(startDate, now time.Time) int {
	days := int(now.Sub(startDate).Hours() / 24)
	if days < 0 {
		return 1
	}
	return days/7 + 1
}
