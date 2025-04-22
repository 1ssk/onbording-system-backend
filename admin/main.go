package main

import (
	"github.com/1ssk/admin-onbording/controllers"
	"github.com/1ssk/admin-onbording/middleware"
	"github.com/1ssk/admin-onbording/initializers"
	"github.com/gin-gonic/gin"
)

func Init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDb()
	initializers.SyncDatabase()
}

func main() {
	r := gin.Default()

	Init()

	adminGroup := r.Group("/", middleware.RequireAdmin) 

	// courseController.go
	adminGroup.POST("/createCourse", controllers.CreateCourse)
	adminGroup.GET("/getCourse", controllers.GetCourse)
	adminGroup.PUT("/updateCourse/:id", controllers.UpdateCourse)
	adminGroup.DELETE("/deleteCourse/:id", controllers.DeleteCourse)

	// lessonController.go
	adminGroup.POST("/createLesson", controllers.CreateLesson)
	adminGroup.GET("/getLesson", controllers.GetLesson)
	adminGroup.PUT("/updateLesson/:id", controllers.UpdateLesson)
	adminGroup.DELETE("/deleteLesson/:id", controllers.DeleteLesson)

	// lessonAttachmentController.go
	adminGroup.POST("/createLessonAttachment", controllers.CreateLessonAttachment)
	adminGroup.GET("/getLessonAttachment", controllers.GetLessonAttachment)
	adminGroup.PUT("/updateLessonAttachment/:id", controllers.UpdateLessonAttachment)
	adminGroup.DELETE("/deleteLessonAttachment/:id", controllers.DeleteLessonAttachment)

	// testController.go
	adminGroup.POST("/createTest", controllers.CreateTest)
	adminGroup.GET("/getTest", controllers.GetTest)
	adminGroup.PUT("/updateTest/:id", controllers.UpdateTest)
	adminGroup.DELETE("/deleteTest/:id", controllers.DeleteTest)

	// courseAndUserController.go
	adminGroup.POST("/createUserAndCourse", controllers.CreateCourseAndUser)
	adminGroup.GET("/getUserAndCourse", controllers.GetCourseAndUser)
	adminGroup.DELETE("/deleteUserAndCourse/:id", controllers.DeleteCourseAndUser)
	
	



	r.Run()
}
