package initializers

import "github.com/1ssk/go-jwt/models"

func SyncDatabase() {
	DB.AutoMigrate(&models.User{})
}
