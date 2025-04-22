package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/1ssk/user-onbording/initializers"
    "github.com/1ssk/user-onbording/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func RequireUser(c *gin.Context) {
    var tokenString string


    authHeader := c.GetHeader("Authorization")
    if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
        tokenString = strings.TrimPrefix(authHeader, "Bearer ")
    } else {

        cookieToken, err := c.Cookie("Authorization")
        if err != nil {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing token"})
            return
        }
        tokenString = cookieToken
    }

    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        secret := os.Getenv("SECRET")
        if secret == "" {
            return nil, fmt.Errorf("secret is not defined")
        }
        return []byte(secret), nil
    })

    if err != nil || !token.Valid {
        c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
        return
    }

    claims, ok := token.Claims.(jwt.MapClaims)
    if !ok {
        c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid claims"})
        return
    }

    if exp, ok := claims["exp"].(float64); ok && int64(exp) < time.Now().Unix() {
        c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token expired"})
        return
    }

    userID, ok := claims["sub"].(float64)
    if !ok {
        c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in token"})
        return
    }

    var user models.User
    result := initializers.DB.First(&user, uint(userID))
    if result.Error != nil {
        c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
        return
    }

    c.Set("user", user)

    c.Next()
}
