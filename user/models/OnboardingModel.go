package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model `gorm:"unique"`
	Email      string `gorm:"unique;not null"`
	Password   string
	Role       string `gorm:"type:varchar(50);not null"`
}

// Course - курс в онбординг-системе
type Course struct {
	gorm.Model
	Title       string `gorm:"type:varchar(255);not null"`
	Description string `gorm:"type:text"`
	Lesson      []Lesson
	UserCourses []UserCourse
}

// UserCourse - связь с пользователями из внешнего сервиса
type UserCourse struct {
	gorm.Model
	UserID   uint `gorm:"not null"`
	CourseID uint `gorm:"not null"`
	Course   Course `gorm:"foreignKey:CourseID"`
}

// Lessons - единица курса (урок, тест и т.д.)
type Lesson struct {
	gorm.Model
	CourseID    uint   `gorm:"not null"`
	Title        string `gorm:"type:varchar(50);not null"`
	Description string `gorm:"type:text"`
	Deadline    *time.Time
	Tests       []Test             `gorm:"foreignKey:LessonID"`
	Attachments []LessonAttachment `gorm:"foreignKey:LessonID"`
}

// LessonsAttachment - вложение для единицы курса
type LessonAttachment struct {
	gorm.Model
	LessonID    uint   `gorm:"not null"`
	Type        string `gorm:"type:varchar(50);not null"`
	URL         string `gorm:"type:varchar(255);not null"`
	Description string `gorm:"type:text"`
}

// Test - тест в единице курса
type Test struct {
	gorm.Model
	LessonID  uint `gorm:"not null"`
	Questions []Question `gorm:"foreignKey:TestID;constraint:OnDelete:CASCADE"`
}

// Question - вопрос в тесте
type Question struct {
	gorm.Model
	TestID  uint `gorm:"not null"`
	Answers []Answer `gorm:"foreignKey:QuestionID;constraint:OnDelete:CASCADE"`
}

// Answer - ответ на вопрос
type Answer struct {
	gorm.Model
	QuestionID uint `gorm:"not null"`
	Text       string
	IsCorrect  bool
}

// проверка теста
type UserAnswer struct {
    gorm.Model
    UserID     uint `gorm:"not null"`
    QuestionID uint `gorm:"not null"`
    AnswerID   uint `gorm:"not null"`
    IsCorrect  bool `gorm:"not null"`
}