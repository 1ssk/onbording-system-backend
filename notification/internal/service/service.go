package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"notification/internal/api"
	"notification/internal/config"
)

type Service struct {
	api.NotificationServiceServer
	config *config.Config
	logger *zap.Logger
}

func New(config *config.Config, logger *zap.Logger) *Service {
	return &Service{
		config: config,
		logger: logger,
	}
}

func (s *Service) SendNotification(ctx context.Context, request *api.SendNotificationRequest) (*api.SendNotificationResponse, error) {
	err := s.SendMessage(s.config, s.logger, request.Mail, request.Subject, request.Text)
	if err != nil {
		return nil, fmt.Errorf("SendNotification: failed to send message: %w", err)
	}

	id := uuid.New().String()

	return &api.SendNotificationResponse{
		Id: id,
	}, nil
}
