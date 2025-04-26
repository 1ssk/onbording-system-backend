package middleware

import (
    "fmt"
    "net/http"
    "os"
    "time"

	"github.com/1ssk/api-gateway/models"
    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt"
)

// Middleware для проверки авторизации
func RequireAuth(c *gin.Context) {
    tokenString, err := c.Cookie("Authorization")
    if err != nil {
        c.AbortWithStatus(http.StatusUnauthorized)
        return
    }

    // Парсим токен
    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return []byte(os.Getenv("SECRET")), nil
    })

    if err != nil || !token.Valid {
        c.AbortWithStatus(http.StatusUnauthorized)
        return
    }

    claims, ok := token.Claims.(jwt.MapClaims)
    if !ok || float64(time.Now().Unix()) > claims["exp"].(float64) {
        c.AbortWithStatus(http.StatusUnauthorized)
        return
    }

    // Проверяем пользователя через go-jwt сервис
    userID := claims["sub"]
    role := claims["role"]

    // Устанавливаем пользователя и роль в контекст
    c.Set("user", models.User{ID: uint(userID.(float64)), Role: role.(string)})
    c.Set("role", role)

    c.Next()
}

// Middleware для проверки админской роли
func RequireAdmin(c *gin.Context) {
    role, exists := c.Get("role")
    if !exists || role != "admin" {
        c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
        return
    }
    c.Next()
}

func RequireLeader(c *gin.Context) {
    role, exists := c.Get("role")
    if !exists || role != "leader" {
        c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Leader access required"})
        return
    }
    c.Next()
}

func RequireManager(c *gin.Context) {
    role, exists := c.Get("role")
    if !exists || role != "manager" {
        c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Leader access required"})
        return
    }
    c.Next()
}