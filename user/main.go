package main

import (
	"github.com/1ssk/user-onbording/controllers"
	"github.com/1ssk/user-onbording/initializers"
	"github.com/1ssk/user-onbording/middleware"
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

	userGroup := r.Group("/", middleware.RequireUser)

    userGroup.GET("/courses", controllers.GetUserCourses)
    userGroup.GET("/course/:course_id/lessons", controllers.GetCourseLessons)
    userGroup.GET("/lesson/:lesson_id/details", controllers.GetLessonDetails)
	userGroup.POST("/answer/submit", controllers.SubmitAnswer)


	r.Run()
}
