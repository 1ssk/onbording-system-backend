package authorization

import (
	"context"
	"crypto/rsa"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type ContextKey string

const Key = ContextKey("user")

type TokenPayload struct {
	ID   string
	Role string
	Exp  int64 // Unix timestamp
}

func MiddlwareJWT(publicKey *rsa.PublicKey) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get token from header
		headerToken := c.GetHeader("Authorization")
		var tokenString string
		_, err := fmt.Sscanf(headerToken, "Bearer %s", &tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		// Parse token
		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
			return publicKey, nil
		})
		// error handling
		if err != nil || !token.Valid || isExpiredToken(claims) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		user := TokenPayload{
			ID:   claims["id"].(string),
			Role: claims["role"].(string),
			Exp:  int64(claims["exp"].(float64)),
		}

		// Set context
		c.Set(string(Key), user)

		// Next
		c.Next()
	}
}

func FromContext(ctx context.Context) (TokenPayload, error) {
	payload, ok := ctx.Value(string(Key)).(TokenPayload)
	if !ok {
		return TokenPayload{}, fmt.Errorf("user not found")
	}
	return payload, nil
}

func isExpiredToken(claims jwt.MapClaims) bool {
	exp, ok := claims["exp"].(int64)
	if !ok {
		return false
	}
	return exp <= time.Now().Unix()
}
