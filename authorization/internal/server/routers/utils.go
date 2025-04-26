package routers

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/glebkasilov/authorization/internal/service"
	"github.com/golang-jwt/jwt/v5"
)

func GetKey(pathToKey string) *rsa.PublicKey {
	data, err := os.ReadFile(pathToKey)
	if err != nil {
		fmt.Printf("failed to read public key: %v", err)
	}

	key, err := jwt.ParseRSAPublicKeyFromPEM(data)
	if err != nil {
		fmt.Printf("failed to parse public key: %v", err)
	}

	return key
}

func HandleError(err error, ctx *gin.Context) {
	if errors.Is(err, service.ErrBadRequest) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}
	if errors.Is(err, service.ErrInternalServer) {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}
	if errors.Is(err, service.ErrNotFound) {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
	}
}
