package storage

import (
	"fmt"

	"github.com/glebkasilov/authorization/internal/config"
	"github.com/glebkasilov/authorization/internal/domain/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Storage struct {
	db *gorm.DB
}

func New() *Storage {
	cfg := config.Config().Storage

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		cfg.Host, cfg.User, cfg.Pass, cfg.DbName, cfg.Port,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&models.User{})

	return &Storage{db: db}
}
