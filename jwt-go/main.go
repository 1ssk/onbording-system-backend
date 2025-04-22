package main

import (
	"github.com/1ssk/go-jwt/controllers"
	"github.com/1ssk/go-jwt/initializers"
	"github.com/1ssk/go-jwt/middleware"
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

	r.POST("/signup", controllers.Signup)
	r.POST("/login", controllers.Login)
	r.GET("/validate", middleware.RequireAuth, controllers.Validate)
	r.PUT("/admin/update-role/:id", middleware.RequireAuth, middleware.RequireAdmin, controllers.UpdateAdmin)

	r.Run()
}
