package main

import (
    "bytes"
    "fmt"
    "io"
    "net/http"

    "github.com/1ssk/api-gateway/middleware"
    "github.com/gin-gonic/gin"
)

func main() {
    r := gin.Default()

    // Эндпоинты для аутентификации (прокси к jwt-go на порт 8000)
    r.POST("/auth/signup", proxyRequest("http://jwt-go:8000/signup"))
    r.POST("/auth/login", proxyRequest("http://jwt-go:8000/login"))
    r.GET("/auth/validate", middleware.RequireAuth, proxyRequest("http://jwt-go:8000/validate"))
    r.PUT("/auth/admin/update-role/:id", middleware.RequireAuth, middleware.RequireAdmin, proxyRequest("http://jwt-go:8000/admin/update-role/:id"))

    // Эндпоинты для администраторов (прокси к admin на порт 3000)
    adminGroup := r.Group("/admin", middleware.RequireAuth, middleware.RequireAdmin)
    {
        // Курсы
        adminGroup.POST("/courses", proxyRequest("http://admin:3000/createCourse"))
        adminGroup.GET("/courses", proxyRequest("http://admin:3000/getCourse"))
        adminGroup.PUT("/courses/:id", proxyRequest("http://admin:3000/updateCourse/:id"))
        adminGroup.DELETE("/courses/:id", proxyRequest("http://admin:3000/deleteCourse/:id"))

        // Уроки
        adminGroup.POST("/lessons", proxyRequest("http://admin:3000/createLesson"))
        adminGroup.GET("/lessons", proxyRequest("http://admin:3000/getLesson"))
        adminGroup.PUT("/lessons/:id", proxyRequest("http://admin:3000/updateLesson/:id"))
        adminGroup.DELETE("/lessons/:id", proxyRequest("http://admin:3000/deleteLesson/:id"))

        // Вложения уроков
        adminGroup.POST("/lesson-attachments", proxyRequest("http://admin:3000/createLessonAttachment"))
        adminGroup.GET("/lesson-attachments", proxyRequest("http://admin:3000/getLessonAttachment"))
        adminGroup.PUT("/lesson-attachments/:id", proxyRequest("http://admin:3000/updateLessonAttachment/:id"))
        adminGroup.DELETE("/lesson-attachments/:id", proxyRequest("http://admin:3000/deleteLessonAttachment/:id"))

        // Тесты
        adminGroup.POST("/tests", proxyRequest("http://admin:3000/createTest"))
        adminGroup.GET("/tests", proxyRequest("http://admin:3000/getTest"))
        adminGroup.PUT("/tests/:id", proxyRequest("http://admin:3000/updateTest/:id"))
        adminGroup.DELETE("/tests/:id", proxyRequest("http://admin:3000/deleteTest/:id"))

        // Курсы и пользователи
        adminGroup.POST("/user-courses", proxyRequest("http://admin:3000/createUserAndCourse"))
        adminGroup.GET("/user-courses", proxyRequest("http://admin:3000/getUserAndCourse"))
        adminGroup.DELETE("/user-courses/:id", proxyRequest("http://admin:3000/deleteUserAndCourse/:id"))
    }

    // Эндпоинты для пользователей (прокси к user на порт 8080)
    userGroup := r.Group("/user", middleware.RequireAuth)
    {
        userGroup.GET("/courses", proxyRequest("http://user:8080/courses"))
        userGroup.GET("/courses/:course_id/lessons", proxyRequest("http://user:8080/course/:course_id/lessons"))
        userGroup.GET("/lessons/:lesson_id/details", proxyRequest("http://user:8080/lesson/:lesson_id/details"))
        userGroup.POST("/answers", proxyRequest("http://user:8080/answer/submit"))
    }

    r.Run(":8080")
}

func proxyRequest(target string) gin.HandlerFunc {
    return func(c *gin.Context) {
        url := target
        for _, param := range c.Params {
            url = fmt.Sprintf("%s/%s", url, param.Value)
        }
        if c.Request.URL.RawQuery != "" {
            url = fmt.Sprintf("%s?%s", url, c.Request.URL.RawQuery)
        }

        body, err := io.ReadAll(c.Request.Body)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read request body"})
            return
        }

        req, err := http.NewRequest(c.Request.Method, url, bytes.NewReader(body))
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
            return
        }

        for key, values := range c.Request.Header {
            for _, value := range values {
                req.Header.Add(key, value)
            }
        }

        client := &http.Client{}
        resp, err := client.Do(req)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to proxy request"})
            return
        }
        defer resp.Body.Close();

        responseBody, err := io.ReadAll(resp.Body)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response body"})
            return
        }

        for key, values := range resp.Header {
            for _, value := range values {
                c.Header(key, value)
            }
        }

        c.Data(resp.StatusCode, resp.Header.Get("Content-Type"), responseBody)
    }
}