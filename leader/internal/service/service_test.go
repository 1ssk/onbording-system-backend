package service

import (
	"context"
	"testing"

	"role-leader/internal/domain/models"
	"role-leader/internal/logger"
	"role-leader/pkg/api"
)

func TestService_New(t *testing.T) {
	storage := &mockStorage{}
	logger := logger.New()

	s := New(storage, logger)

	if s.storage != storage {
		t.Errorf("New() storage = %v, want %v", s.storage, storage)
	}
	if s.logger != logger {
		t.Errorf("New() logger = %v, want %v", s.logger, logger)
	}
}

func TestService_UpdateMeetingFeedback(t *testing.T) {
	storage := &mockStorage{}
	logger := logger.New()
	s := New(storage, logger)
	req := &api.CreateFeedbackRequest{
		MeetingId: "123",
		Message:   "Test feedback",
	}

	_, err := s.UpdateMeetingFeedback(context.Background(), req)

	if err != nil {
		t.Errorf("UpdateMeetingFeedback() error = %v", err)
	}
}

func TestService_Meeting(t *testing.T) {
	storage := &mockStorage{}
	logger := logger.New()
	s := New(storage, logger)
	req := &api.GetMeetingRequest{
		MeetingId: "123",
	}

	_, err := s.Meeting(context.Background(), req)

	if err != nil {
		t.Errorf("Meeting() error = %v", err)
	}
}

func TestService_Meetings(t *testing.T) {
	storage := &mockStorage{}
	logger := logger.New()
	s := New(storage, logger)
	req := &api.GetLeaderMeetingsRequest{
		LeaderId: "123",
	}

	_, err := s.Meetings(context.Background(), req)

	if err != nil {
		t.Errorf("Meetings() error = %v", err)
	}
}

type mockStorage struct{}

func (m *mockStorage) Meeting(ctx context.Context, id string) (*models.Meeting, error) {
	return &models.Meeting{}, nil
}

func (m *mockStorage) Meetings(ctx context.Context, leaderID string) ([]*models.Meeting, error) {
	return []*models.Meeting{}, nil
}

func (m *mockStorage) UpdateMeetingFeedback(ctx context.Context, id string, feedback string) error {
	return nil
}
