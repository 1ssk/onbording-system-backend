package storage

import (
	"context"

	"github.com/glebkasilov/authorization/internal/domain/models"
)

func (s *Storage) SaveUser(ctx context.Context, user *models.User) error {
	return s.db.Create(&user).Error
}

func (s *Storage) User(ctx context.Context, id string) (*models.User, error) {
	var user models.User

	if err := s.db.First(&user, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *Storage) UserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User

	if err := s.db.First(&user, "email = ?", email).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *Storage) SetRole(ctx context.Context, id string, role string) error {
	return s.db.Model(&models.User{}).Where("id = ?", id).Update("role", role).Error
}
