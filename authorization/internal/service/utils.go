package service

import (
	"time"

	"github.com/glebkasilov/authorization/internal/config"
	"github.com/glebkasilov/authorization/internal/domain/models"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func CreateToken(user *models.User) string {
	exp := time.Now().Add(config.Config().JWT.TokenTTl).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"id":   user.ID,
		"role": user.Role,
		"exp":  exp,
	})

	tokenString, err := token.SignedString(config.Config().JWT.PrivateKey)

	if err != nil {
		return ""
	}

	return tokenString
}
