package storage

import (
	"context"
	"fmt"
	"role-leader/internal/config"
	"role-leader/internal/domain/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Storage struct {
	db *gorm.DB
}

func New(dbcfg config.DatabaseCfg) (*Storage, error) {
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
	var meeting models.Meeting
	err := s.db.WithContext(ctx).
		Model(&models.Meeting{}).
		Where("id = ?", id).
		First(&meeting).
		Error
	return &meeting, err
}

func (s *Storage) Meetings(ctx context.Context, leaderID string) ([]*models.Meeting, error) {
	var meetings []*models.Meeting
	err := s.db.WithContext(ctx).
		Model(&models.Meeting{}).
		Where("leader_id = ?", leaderID).
		Find(&meetings).
		Error
	return meetings, err
}

func (s *Storage) UpdateMeetingFeedback(ctx context.Context, id string, feedback string) error {
	return s.db.WithContext(ctx).
		Model(&models.Meeting{}).
		Where("id = ?", id).
		Update("feedback", feedback).
		Error
}
