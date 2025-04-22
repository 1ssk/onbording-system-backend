package initializers

import "github.com/1ssk/admin-onbording/models"

func SyncDatabase() {
	DB.AutoMigrate(
		&models.User{},
		&models.Course{},
		&models.Lesson{},
		&models.LessonAttachment{},
		&models.Test{},
		&models.Question{},
		&models.Answer{},
		&models.UserCourse{},
		&models.UserAnswer{},
	)
}
