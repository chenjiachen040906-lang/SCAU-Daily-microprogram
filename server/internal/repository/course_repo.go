package repository

import (
	"scau-daily/internal/database"
	"scau-daily/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func GetCurrentSemester() (*model.Semester, error) {
	var semester model.Semester
	if err := database.DB.Where("is_current = ?", true).First(&semester).Error; err != nil {
		return nil, err
	}
	return &semester, nil
}

func GetCoursesByUserAndSemester(userID uuid.UUID, semesterID uint) ([]model.Course, error) {
	var courses []model.Course
	err := database.DB.
		Where("user_id = ? AND semester_id = ?", userID, semesterID).
		Preload("Schedules").
		Find(&courses).Error
	return courses, err
}

func CreateCourse(course *model.Course) error {
	return database.DB.Create(course).Error
}

func BatchCreateCourses(courses []model.Course) error {
	if len(courses) == 0 {
		return nil
	}
	return database.DB.Create(&courses).Error
}

func DeleteCoursesByUserAndSemester(userID uuid.UUID, semesterID uint) error {
	return database.DB.Transaction(func(tx *gorm.DB) error {
		// Collect course IDs first
		var courseIDs []uuid.UUID
		if err := tx.Model(&model.Course{}).
			Where("user_id = ? AND semester_id = ?", userID, semesterID).
			Pluck("id", &courseIDs).Error; err != nil {
			return err
		}
		if len(courseIDs) == 0 {
			return nil
		}

		// Delete child schedules
		if err := tx.Where("course_id IN ?", courseIDs).Delete(&model.CourseSchedule{}).Error; err != nil {
			return err
		}

		// Delete courses
		if err := tx.Where("id IN ?", courseIDs).Delete(&model.Course{}).Error; err != nil {
			return err
		}

		return nil
	})
}
