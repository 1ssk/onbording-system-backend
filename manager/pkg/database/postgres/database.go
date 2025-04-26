package postgres

import (
	"context"
	"fmt"
	"strconv"

	"github.com/glebkasilov/grpc-manager/internal/config"
	"github.com/glebkasilov/grpc-manager/internal/domain/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Storage struct {
	db *gorm.DB
}

func New(dbcfg *config.DatabaseCfg) (*Storage, error) {

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d",
		dbcfg.Host,
		dbcfg.User,
		dbcfg.Password,
		dbcfg.DbName,
		dbcfg.Port,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(&models.User{}, &models.Meeting{}); err != nil {
		return nil, fmt.Errorf("failed to auto migrate: %w", err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) Close() error {
	db, err := s.db.DB()
	if err != nil {
		return err
	}
	return db.Close()
}

func (s *Storage) Meeting(ctx context.Context, id string) (*models.Meeting, error) {
	uintID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return nil, err
	}

	var meeting models.Meeting
	err = s.db.WithContext(ctx).
		Model(&models.Meeting{}).
		Where("id = ?", uintID).
		First(&meeting).
		Error
	return &meeting, err
}

func (s *Storage) AddMeeting(ctx context.Context, meeting *models.Meeting) error {
	return s.db.WithContext(ctx).Create(meeting).Error
}

func (s *Storage) Meetings(ctx context.Context) ([]models.Meeting, error) {
	var meetings []models.Meeting
	err := s.db.WithContext(ctx).Find(&meetings).Error
	return meetings, err
}

func (s *Storage) User(ctx context.Context, id string) (*models.User, error) {
	var user models.User
	err := s.db.WithContext(ctx).
		Where("id = ?", id).
		First(&user).
		Error
	return &user, err
}

func (s *Storage) UserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := s.db.WithContext(ctx).
		Where("email = ?", email).
		First(&user).
		Error
	return &user, err
}

func (s *Storage) CompleteCourse(ctx context.Context, email string) error {
	var user models.User
	if err := s.db.WithContext(ctx).
		Where("email = ?", email).
		First(&user).
		Error; err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	return s.db.WithContext(ctx).
		Model(&user).
		Update("current_stage", "completed").
		Error
}

func (s *Storage) UpdateMeeting(ctx context.Context, meeting *models.Meeting) error {
	return s.db.WithContext(ctx).
		Model(meeting).
		Updates(meeting).
		Error
}

func (s *Storage) DeleteMeeting(ctx context.Context, id string) error {
	return s.db.WithContext(ctx).
		Where("id = ?", id).
		Delete(&models.Meeting{}).
		Error
}

func (s *Storage) AddUser(ctx context.Context, user *models.User) error {
	return s.db.WithContext(ctx).Create(user).Error
}
