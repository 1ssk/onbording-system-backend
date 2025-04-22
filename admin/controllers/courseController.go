package controllers

import (
	"github.com/1ssk/admin-onbording/initializers"
	"github.com/1ssk/admin-onbording/models"
	"github.com/gin-gonic/gin"
)

func CreateCourse(c *gin.Context) {
	var course struct {
		Title       string
		Description string
	}

	c.Bind(&course)

	newCourse := models.Course{Title: course.Title, Description: course.Description}
	result := initializers.DB.Create(&newCourse)

	if result.Error != nil {
		c.Status(400)
		return
	}

	c.JSON(200, gin.H{
		"newCourse": course,
	})
}

func GetCourse(c *gin.Context) {

	var course []models.Course
	initializers.DB.Find(&course)

	c.JSON(200, gin.H{
		"course": course,
	})
}

func UpdateCourse(c *gin.Context) {
	id := c.Param("id")

	var course struct {
		Title       string
		Description string
	}

	c.Bind(&course)

	var updateCourse models.Course
	initializers.DB.First(&updateCourse, id)

	initializers.DB.Model(&updateCourse).Updates(models.Course{
		Title:       course.Title,
		Description: course.Description,
	})

	c.JSON(200, gin.H{
		"course": course,
	})

}

func DeleteCourse(c *gin.Context) {
	id := c.Param("id")

	var course models.Course
	initializers.DB.Delete(&course, id)

	c.JSON(200, gin.H{
		"course": "acces delete",
	})
}
