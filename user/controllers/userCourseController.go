package controllers

import (
	"github.com/1ssk/user-onbording/initializers"
    "github.com/1ssk/user-onbording/models"
    "github.com/gin-gonic/gin"
	"strconv"
)



func GetUserCourses(c *gin.Context) {
    user, exists := c.Get("user")
    if !exists {
        c.JSON(401, gin.H{"error": "unauthorized"})
        return
    }

    u, ok := user.(models.User)
    if !ok {
        c.JSON(401, gin.H{"error": "invalid user data"})
        return
    }

    var userCourses []models.UserCourse
    result := initializers.DB.Preload("Course").Where("user_id = ?", u.ID).Find(&userCourses)
    if result.Error != nil {
        c.JSON(500, gin.H{"error": result.Error.Error()})
        return
    }

    var courses []models.Course
    for _, uc := range userCourses {
        courses = append(courses, uc.Course)
    }

    c.JSON(200, gin.H{"courses": courses})
}

func GetCourseLessons(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(401, gin.H{"error": "unauthorized"})
		return
	}
	u := user.(models.User)

	courseIDStr := c.Param("course_id")
	courseID, err := strconv.ParseUint(courseIDStr, 10, 32)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid course_id"})
		return
	}

	var userCourse models.UserCourse
	result := initializers.DB.Where("user_id = ? AND course_id = ?", u.ID, courseID).First(&userCourse)
	if result.Error != nil {
		c.JSON(403, gin.H{"error": "course not accessible for this user"})
		return
	}

	var lessons []models.Lesson
	result = initializers.DB.Where("course_id = ?", courseID).Find(&lessons)
	if result.Error != nil {
		c.JSON(500, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(200, gin.H{"lessons": lessons})
}

func GetLessonDetails(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(401, gin.H{"error": "unauthorized"})
		return
	}
	u := user.(models.User)

	lessonIDStr := c.Param("lesson_id")
	lessonID, err := strconv.ParseUint(lessonIDStr, 10, 32)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid lesson_id"})
		return
	}

	var lesson models.Lesson
	result := initializers.DB.
		Preload("Attachments").
		Preload("Tests.Questions.Answers").
		First(&lesson, lessonID)
	if result.Error != nil {
		c.JSON(404, gin.H{"error": "lesson not found"})
		return
	}

	var userCourse models.UserCourse
	result = initializers.DB.Where("user_id = ? AND course_id = ?", u.ID, lesson.CourseID).First(&userCourse)
	if result.Error != nil {
		c.JSON(403, gin.H{"error": "lesson not accessible for this user"})
		return
	}

	c.JSON(200, gin.H{
		"lesson": gin.H{
			"id":          lesson.ID,
			"title":       lesson.Title,
			"description": lesson.Description,
			"deadline":    lesson.Deadline,
			"attachments": lesson.Attachments,
			"tests":       lesson.Tests,
		},
	})
}

func SubmitAnswer(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(401, gin.H{"error": "unauthorized"})
		return
	}
	u := user.(models.User)

	var req struct {
		QuestionID uint `json:"question_id"`
		AnswerID   uint `json:"answer_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "invalid request body"})
		return
	}

	var question models.Question
	if result := initializers.DB.First(&question, req.QuestionID); result.Error != nil {
		c.JSON(404, gin.H{"error": "question not found"})
		return
	}

	var answer models.Answer
	if result := initializers.DB.Where("id = ? AND question_id = ?", req.AnswerID, req.QuestionID).First(&answer); result.Error != nil {
		c.JSON(404, gin.H{"error": "answer not found or does not belong to this question"})
		return
	}

	var test models.Test
	if result := initializers.DB.First(&test, question.TestID); result.Error != nil {
		c.JSON(404, gin.H{"error": "test not found"})
		return
	}

	var lesson models.Lesson
	if result := initializers.DB.First(&lesson, test.LessonID); result.Error != nil {
		c.JSON(404, gin.H{"error": "lesson not found"})
		return
	}

	var userCourse models.UserCourse
	if result := initializers.DB.Where("user_id = ? AND course_id = ?", u.ID, lesson.CourseID).First(&userCourse); result.Error != nil {
		c.JSON(403, gin.H{"error": "test not accessible for this user"})
		return
	}

	userAnswer := models.UserAnswer{
		UserID:     u.ID,
		QuestionID: req.QuestionID,
		AnswerID:   req.AnswerID,
		IsCorrect:  answer.IsCorrect,
	}

	if result := initializers.DB.Create(&userAnswer); result.Error != nil {
		c.JSON(500, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(200, gin.H{
		"message":     "answer submitted",
		"is_correct":  answer.IsCorrect,
		"question_id": req.QuestionID,
		"answer_id":   req.AnswerID,
	})
}
