package service

import (
	"context"
	"errors"
	"time"

	"github.com/glebkasilov/grpc-manager/internal/config"
	"github.com/glebkasilov/grpc-manager/internal/domain/models"
	"github.com/redis/go-redis/v9"

	test "github.com/glebkasilov/grpc-manager/pkg/api"
	"github.com/glebkasilov/grpc-manager/pkg/database/postgres"
	"github.com/glebkasilov/grpc-manager/pkg/logger"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type Storage interface {
	CompleteCourse(ctx context.Context, email string) error
	UserByEmail(ctx context.Context, email string) (*models.User, error)
	AddUser(ctx context.Context, user *models.User) error
	AddMeeting(ctx context.Context, meeting *models.Meeting) error
	Meetings(ctx context.Context) ([]models.Meeting, error)
	Meeting(ctx context.Context, id string) (*models.Meeting, error)
	UpdateMeeting(ctx context.Context, meeting *models.Meeting) error
	DeleteMeeting(ctx context.Context, id string) error
}

type Service struct {
	test.MeetingServiceServer
	storage postgres.Storage
	redis   *redis.Client
	logger  *logger.Logger
	ttl     time.Duration
}

func New(storage postgres.Storage, redisClient *redis.Client, logger *logger.Logger) *Service {
	return &Service{
		storage: storage,
		redis:   redisClient,
		logger:  logger,
		ttl:     config.Config().Redis.TTL,
	}
}

func (s *Service) AddUser(ctx context.Context, req *test.AddUserRequest) (*test.AddUserResponse, error) {
	s.logger.Info("AddUser called", zap.String("email", req.Email))

	user := &models.User{
		Email:        req.Email,
		Password:     req.Password,
		Role:         "user",
	}

	if err := s.storage.AddUser(ctx, user); err != nil {
		s.logger.Error("Failed to add user", zap.Error(err))
		return nil, err
	}

	s.logger.Info("User added", zap.Uint("id", user.ID))
	return &test.AddUserResponse{
		Id: string(user.ID),
	}, nil
}

func (s *Service) AddLeader(ctx context.Context, req *test.AddLeaderRequest) (*test.AddLeaderResponse, error) {
	s.logger.Info("AddLeader called", zap.String("email", req.Email))

	user := &models.User{
		Email:        req.Email,
		Password:     req.Password,
		Role:         "leader",
	}

	if err := s.storage.AddUser(ctx, user); err != nil {
		s.logger.Error("Failed to add leader", zap.Error(err))
		return nil, err
	}

	s.logger.Info("Leader added", zap.Uint("id", user.ID))
	return &test.AddLeaderResponse{
		Id: string(user.ID),
	}, nil
}

func (s *Service) FinishCourse(ctx context.Context, req *test.FinishCourseRequest) (*test.FinishCourseResponse, error) {
	s.logger.Info("CompleteCourse called", zap.String("email", req.Email))

	if err := s.storage.CompleteCourse(ctx, req.Email); err != nil {
		s.logger.Error("Failed to complete course", zap.Error(err))
		return nil, err
	}

	s.logger.Info("Course completed", zap.String("email", req.Email))
	return &test.FinishCourseResponse{
		Id: req.Email,
	}, nil
}

func (s *Service) AddMeeting(ctx context.Context, req *test.AddMeetingRequest) (*test.AddMeetingResponse, error) {
	s.logger.Info("AddMeeting called")

	time, err := convertTime(req.StartTime)
	if err != nil {
		s.logger.Error("Failed to convert time", zap.Error(err))
		return nil, err
	}

	meeting := &models.Meeting{
		UserID:    req.UserId,
		LeaderID:  req.LeaderId,
		Title:     req.Title,
		StartTime: &time,
		Status:    "pending",
		Feedback:  "",
	}

	if err := s.storage.AddMeeting(ctx, meeting); err != nil {
		s.logger.Error("Failed to add meeting", zap.Error(err))
		return nil, err
	}

	s.logger.Info("Meeting added", zap.Uint("id", meeting.ID))
	return &test.AddMeetingResponse{
		Id: string(meeting.ID),
	}, nil
}

func (s *Service) GetMeetings(ctx context.Context, req *test.GetMeetingsRequest) (*test.GetMeetingsResponse, error) {
	s.logger.Info("GetMeetings called",
		zap.String("caller", getCallerInfo(ctx)))

	meetings, err := s.storage.Meetings(ctx)
	if err != nil {
		s.logger.Error("Failed to fetch meetings",
			zap.Error(err),
			zap.String("operation", "get_meetings"))
		return nil, status.Error(codes.Internal, "internal server error")
	}

	s.logger.Debug("Successfully fetched meetings",
		zap.Int("count", len(meetings)),
		zap.Any("first_meeting_id", zap.Uint("id", meetings[0].ID)))

	pbMeetings := make([]*test.Meeting, 0, len(meetings))
	for _, m := range meetings {
		startTime := m.StartTime.String()

		pbMeeting := &test.Meeting{
			Id:        string(m.ID),
			UserId:    m.UserID,
			LeaderId:  m.LeaderID,
			Title:     m.Title,
			StartTime: startTime, // UTC + наносекунды
		}

		pbMeetings = append(pbMeetings, pbMeeting)
	}

	s.logger.Debug("Successfully retrieved meetings",
		zap.Int("count", len(pbMeetings)),
		zap.Any("first_meeting", safeLogMeeting(pbMeetings[0])))

	return &test.GetMeetingsResponse{Meetings: pbMeetings}, nil
}

func (s *Service) GetMeeting(ctx context.Context, req *test.GetMeetingRequest) (*test.GetMeetingResponse, error) {
	s.logger.Info("GetMeeting called",
		zap.String("meeting_id", req.MeetingId),
		zap.String("caller", getCallerInfo(ctx)),
	)

	if req.MeetingId == "" {
		return nil, status.Error(codes.InvalidArgument, "meeting ID is required")
	}

	meeting, err := s.storage.Meeting(ctx, req.MeetingId)
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		s.logger.Debug("Meeting not found",
			zap.String("meeting_id", req.MeetingId))
		return nil, status.Error(codes.NotFound, "meeting not found")
	case err != nil:
		s.logger.Error("Database error",
			zap.Error(err),
			zap.String("meeting_id", req.MeetingId))
		return nil, status.Error(codes.Internal, "database error")
	}

	startTimeUTC := meeting.StartTime.UTC()

	pbMeeting := &test.Meeting{
		Id:        string(meeting.ID),
		UserId:    meeting.UserID,
		LeaderId:  meeting.LeaderID,
		Title:     meeting.Title,
		StartTime: startTimeUTC.Format(time.RFC3339Nano),
	}

	s.logger.Debug("Meeting details retrieved",
		zap.String("meeting_id", pbMeeting.Id),
		zap.Time("start_time", startTimeUTC))

	return &test.GetMeetingResponse{Meeting: pbMeeting}, nil
}

func (s *Service) UpdateMeeting(ctx context.Context, req *test.UpdateMeetingRequest) (*test.UpdateMeetingResponse, error) {
	s.logger.Info("UpdateMeeting called", zap.String("meeting_id", req.MeetingId))

	meeting, err := s.storage.Meeting(ctx, req.MeetingId)
	if err != nil {
		s.logger.Error("Failed to get meeting for update", zap.Error(err))
		return nil, err
	}

	meeting.UserID = req.UserId
	meeting.LeaderID = req.LeaderId
	meeting.Title = req.Title

	s.logger.Info("Meeting updated", zap.Uint("id", meeting.ID))

	if err := s.storage.UpdateMeeting(ctx, meeting); err != nil {
		s.logger.Error("Failed to update meeting", zap.Error(err))
		return nil, err
	}

	s.logger.Info("Meeting updated", zap.Uint("id", meeting.ID))
	return &test.UpdateMeetingResponse{}, nil
}

func (s *Service) DeleteMeeting(ctx context.Context, req *test.DeleteMeetingRequest) (*test.DeleteMeetingResponse, error) {
	s.logger.Info("DeleteMeeting called", zap.String("meeting_id", req.MeetingId))

	if err := s.storage.DeleteMeeting(ctx, req.MeetingId); err != nil {
		s.logger.Error("Failed to delete meeting", zap.Error(err))
		return nil, err
	}

	s.logger.Info("Meeting deleted", zap.String("id", req.MeetingId))
	return &test.DeleteMeetingResponse{}, nil
}

func (s *Service) CreateUser(ctx context.Context, req *test.AddUserRequest) (*test.AddUserResponse, error) {
	s.logger.Info("CreateUser called", zap.String("email", req.Email))

	id := uuid.New().String()

	hashedPassword, err := HashPassword(req.Password)
	if err != nil {
		s.logger.Error("Failed to hash password", zap.Error(err))
		return nil, err
	}

	user := &models.User{
		Email:        req.Email,
		Password:     hashedPassword,
		Role:         "user",
	}
	if err := s.storage.AddUser(ctx, user); err != nil {
		s.logger.Error("Failed to create user", zap.Error(err))
		return nil, err
	}

	s.logger.Info("User created", zap.String("id", id))
	return &test.AddUserResponse{Id: id}, nil

}

func (s *Service) GetUser(ctx context.Context, req *test.GetUserRequest) (*test.GetUserResponse, error) {
	s.logger.Info("UserByEmail called",
		zap.String("email", req.Email),
		zap.String("caller", getCallerInfo(ctx)),
	)

	user, err := s.storage.UserByEmail(ctx, req.Email)
	if err != nil {
		s.logger.Error("Failed to get user by email", zap.Error(err))
		return nil, err
	}

	s.logger.Info("User found", zap.Uint("id", user.ID))
	return &test.GetUserResponse{
		Id:           string(user.ID),
		Email:        user.Email,
	}, nil
}
