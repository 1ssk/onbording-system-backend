package main

import (
    "bytes"
    "io"
    "net/http"
    "strings"

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

    // Группа эндпоинтов для менеджера
    managerGroup := r.Group("/manager", middleware.RequireAuth, middleware.RequireManager)
    {
        managerGroup.POST("/meetings", proxyRequest("http://manager:8082/api/meetings"))                                     // AddMeeting
        managerGroup.GET("/meetings/:meeting_id", proxyRequest("http://manager:8082/api/meetings/:meeting_id"))             // GetMeeting
        managerGroup.PATCH("/meetings/:meeting_id", proxyRequest("http://manager:8082/api/meetings/:meeting_id"))           // UpdateMeeting
        managerGroup.DELETE("/meetings/:meeting_id", proxyRequest("http://manager:8082/api/meetings/:meeting_id"))          // DeleteMeeting
        managerGroup.GET("/meetings", proxyRequest("http://manager:8082/api/meetings"))                                      // GetMeetings
        managerGroup.POST("/meetings/users", proxyRequest("http://manager:8082/api/meetings/users"))                         // AddUser
        managerGroup.GET("/meetings/users/:email", proxyRequest("http://manager:8082/api/meetings/users/:email"))           // GetUser
        managerGroup.POST("/meetings/leaders", proxyRequest("http://manager:8082/api/meetings/leaders"))                     // AddLeader
        managerGroup.POST("/meetings/users/:email/finish", proxyRequest("http://manager:8082/api/meetings/users/:email/finish")) // FinishCourse
    }

    leaderGroup := r.Group("/leader", middleware.RequireAuth, middleware.RequireLeader)
    {
        leaderGroup.POST("/feedback", proxyRequest("http://leader:8083/api/feedback"))                                                  // CreateFeedback
        leaderGroup.GET("/meetings/:meeting_id", proxyRequest("http://leader:8083/api/leaders/meetings/:meeting_id"))                 // GetMeeting
        leaderGroup.GET("/meetings-by-leader/:leader_id", proxyRequest("http://leader:8083/api/leaders/:leader_id/meetings"))         // GetLeaderMeetings
    }

    notificationGroup := r.Group("/api", middleware.RequireAuth, middleware.RequireAdmin)
    {
        notificationGroup.POST("/notification", proxyRequest("http://notification:8084/api/notifications"))                                                
    }



    r.Run(":8080")
}

func proxyRequest(target string) gin.HandlerFunc {
    return func(c *gin.Context) {
        // Подставляем параметры из маршрута (например :id -> 123)
        targetURL := target
        for _, param := range c.Params {
            placeholder := ":" + param.Key
            targetURL = strings.ReplaceAll(targetURL, placeholder, param.Value)
        }

        // Добавляем query string если есть
        if c.Request.URL.RawQuery != "" {
            targetURL += "?" + c.Request.URL.RawQuery
        }

        // Читаем тело запроса
        body, err := io.ReadAll(c.Request.Body)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось прочитать тело запроса"})
            return
        }

        // Создаем новый запрос на целевой URL
        req, err := http.NewRequest(c.Request.Method, targetURL, bytes.NewReader(body))
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось создать новый запрос"})
            return
        }

        // Копируем заголовки
        for key, values := range c.Request.Header {
            for _, value := range values {
                req.Header.Add(key, value)
            }
        }

        // Отправляем проксированный запрос
        client := &http.Client{}
        resp, err := client.Do(req)
        if err != nil {
            c.JSON(http.StatusBadGateway, gin.H{"error": "Ошибка проксирования запроса"})
            return
        }
        defer resp.Body.Close()

        // Читаем ответ от проксируемого сервиса
        responseBody, err := io.ReadAll(resp.Body)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось прочитать ответ от сервиса"})
            return
        }

        // Передаем заголовки ответа клиенту
        for key, values := range resp.Header {
            for _, value := range values {
                c.Writer.Header().Add(key, value)
            }
        }

        // Отдаем ответ клиенту
        c.Data(resp.StatusCode, resp.Header.Get("Content-Type"), responseBody)
    }
}
