package controllers

import (
	"github.com/1ssk/admin-onbording/initializers"
	"github.com/1ssk/admin-onbording/models"
	"github.com/gin-gonic/gin"
)

func CreateCourseAndUser(c *gin.Context) {

	var userAndCourse struct {
		UserID   uint 
		CourseID uint 
	}

	c.Bind(&userAndCourse)

	newUserAndCourse := models.UserCourse{UserID: userAndCourse.UserID, CourseID: userAndCourse.CourseID}
	result := initializers.DB.Create(&newUserAndCourse)

	if result.Error != nil {
		c.Status(400)
		return
	}

	c.JSON(200, gin.H{
		"UserAndCourse": userAndCourse,
	})

}

func GetCourseAndUser(c *gin.Context) {

	var userAndCourse []models.UserCourse
	initializers.DB.Find(&userAndCourse)

	c.JSON(200, gin.H{
		"userAndCourse": userAndCourse,
	})
}


func DeleteCourseAndUser(c *gin.Context) {
	id := c.Param("id")

	var userAndCourse models.UserCourse
	initializers.DB.Delete(&userAndCourse, id)

	c.JSON(200, gin.H{
		"userAndCourse": "acces delete",
	})
}
