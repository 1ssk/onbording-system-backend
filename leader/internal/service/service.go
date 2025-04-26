package service

import (
	"context"
	"role-leader/internal/domain/models"
	"role-leader/internal/logger"
	"role-leader/pkg/api"

	"strconv"
	"go.uber.org/zap"
)

type Storage interface {
	Meeting(ctx context.Context, id string) (*models.Meeting, error)
	Meetings(ctx context.Context, leaderID string) ([]*models.Meeting, error)
	UpdateMeetingFeedback(ctx context.Context, id string, feedback string) error
}

type Service struct {
	api.RoleLeaderServer
	storage Storage
	logger  *logger.Logger
}

func New(storage Storage, logger *logger.Logger) *Service {
	return &Service{
		storage: storage,
		logger:  logger,
	}
}

func (s *Service) UpdateMeetingFeedback(ctx context.Context, req *api.CreateFeedbackRequest) (*api.CreateFeedbackResponse, error) {
	s.logger.Info("CreateFeedback", zap.String("call_id", req.MeetingId), zap.String("message", req.Message))

	if err := s.storage.UpdateMeetingFeedback(ctx, req.MeetingId, req.Message); err != nil {
		return nil, err
	}

	return &api.CreateFeedbackResponse{Status: "ok"}, nil
}

func (s *Service) Meeting(ctx context.Context, req *api.GetMeetingRequest) (*api.GetMeetingResponse, error) {
	s.logger.Info("GetMeeting", zap.String("call_id", req.MeetingId))

	meeting, err := s.storage.Meeting(ctx, req.MeetingId)
	if err != nil {
		return nil, err
	}

	meetingAns := &api.Meeting{
		MeetingId: strconv.FormatUint(uint64(meeting.ID), 10),
		UserId:    meeting.UserID,
		LeaderId:  meeting.LeaderID,
		Title:     meeting.Title,
		StartTime: meeting.StartTime.String(),
		Status:    meeting.Status,
		Feedback:  meeting.Feedback,
	}

	return &api.GetMeetingResponse{Meeting: meetingAns}, nil
}

func (s *Service) Meetings(ctx context.Context, req *api.GetLeaderMeetingsRequest) (*api.GetLeaderMeetingsResponse, error) {
	s.logger.Info("GetLeaderMeetings", zap.String("leader_id", req.LeaderId))

	meetings, err := s.storage.Meetings(ctx, req.LeaderId)
	if err != nil {
		return nil, err
	}

	var meetingsAns []*api.Meeting
	for _, meeting := range meetings {
		meetingAns := &api.Meeting{
			MeetingId: strconv.FormatUint(uint64(meeting.ID), 10),
			UserId:    meeting.UserID,
			LeaderId:  meeting.LeaderID,
			Title:     meeting.Title,
			StartTime: meeting.StartTime.String(),
			Status:    meeting.Status,
			Feedback:  meeting.Feedback,
		}
		meetingsAns = append(meetingsAns, meetingAns)
	}

	return &api.GetLeaderMeetingsResponse{Meetings: meetingsAns}, nil
}
