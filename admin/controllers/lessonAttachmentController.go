package controllers

import (
	"github.com/1ssk/admin-onbording/initializers"
	"github.com/1ssk/admin-onbording/models"
	"github.com/gin-gonic/gin"
)

func CreateLessonAttachment(c *gin.Context) {

	var lessonAttachment struct {
		LessonID    uint  
		Type        string 
		URL         string 
		Description string 
	}

	c.Bind(&lessonAttachment)

	newLessonAttachment := models.LessonAttachment{
		LessonID: lessonAttachment.LessonID, 
		Type: lessonAttachment.Type, 
		URL: lessonAttachment.URL, 
		Description: lessonAttachment.Description,
	}
	
	result := initializers.DB.Create(&newLessonAttachment)

	if result.Error != nil {
		c.Status(400)
		return
	}

	c.JSON(200, gin.H{
		"lessonAttachment": lessonAttachment,
	})
}

func GetLessonAttachment(c *gin.Context) {

	var lessonAttachment []models.LessonAttachment
	initializers.DB.Find(&lessonAttachment)

	c.JSON(200, gin.H{
		"lessonAttachment": lessonAttachment,
	})
}

func UpdateLessonAttachment(c *gin.Context) {
	id := c.Param("id")

	var lessonAttachment struct {
		LessonID    uint  
		Type        string 
		URL         string 
		Description string 
	}


	c.Bind(&lessonAttachment)

	var updateLessonAttachment models.LessonAttachment
	initializers.DB.First(&updateLessonAttachment, id)

	initializers.DB.Model(&updateLessonAttachment).Updates(models.LessonAttachment{
		LessonID: lessonAttachment.LessonID, 
		Type: lessonAttachment.Type, 
		URL: lessonAttachment.URL, 
		Description: lessonAttachment.Description,
	})

	c.JSON(200, gin.H{
		"lessonAttachment": lessonAttachment,
	})

}

func DeleteLessonAttachment(c *gin.Context) {
	id := c.Param("id")

	var lessonAttachment models.LessonAttachment
	initializers.DB.Delete(&lessonAttachment, id)

	c.JSON(200, gin.H{
		"lessonAttachment": "acces delete",
	})
}
