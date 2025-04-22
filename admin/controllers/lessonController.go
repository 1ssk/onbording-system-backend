package controllers

import (
	"github.com/1ssk/admin-onbording/initializers"
	"github.com/1ssk/admin-onbording/models"
	"github.com/gin-gonic/gin"
	"time"
)

func CreateLesson(c *gin.Context) {

	var lesson struct {
		CourseID    uint   
		Title       string 
		Description string
		Deadline    *time.Time
	}

	c.Bind(&lesson)

	newLesson := models.Lesson{CourseID: lesson.CourseID, Title: lesson.Title, Description: lesson.Description, Deadline: lesson.Deadline}
	result := initializers.DB.Create(&newLesson)

	if result.Error != nil {
		c.Status(400)
		return
	}

	c.JSON(200, gin.H{
		"lesson": lesson,
	})
}

func GetLesson(c *gin.Context) {

	var lesson []models.Lesson
	initializers.DB.Find(&lesson)

	c.JSON(200, gin.H{
		"lesson": lesson,
	})
}

func UpdateLesson(c *gin.Context) {
	id := c.Param("id")

	var lesson struct {
		CourseID    uint   
		Title       string 
		Description string
		Deadline    *time.Time
	}


	c.Bind(&lesson)

	var updateLesson models.Lesson
	initializers.DB.First(&updateLesson, id)

	initializers.DB.Model(&updateLesson).Updates(models.Lesson{
		CourseID: lesson.CourseID, 
		Title: lesson.Title, 
		Description: lesson.Description, 
		Deadline: lesson.Deadline,
	})

	c.JSON(200, gin.H{
		"lesson": lesson,
	})

}

func DeleteLesson(c *gin.Context) {
	id := c.Param("id")

	var lesson models.Lesson
	initializers.DB.Delete(&lesson, id)

	c.JSON(200, gin.H{
		"lesson": "acces delete",
	})
}
