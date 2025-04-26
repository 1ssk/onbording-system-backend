package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/glebkasilov/authorization/internal/domain/models"
	"github.com/glebkasilov/authorization/internal/domain/requests"
	"github.com/google/uuid"
)

var (
	ErrInternalServer = errors.New("internal server error")
	ErrNotFound       = errors.New("not found")
	ErrBadRequest     = errors.New("bad request")
)

type Repository interface {
	SaveUser(ctx context.Context, user *models.User) error
	User(ctx context.Context, id string) (*models.User, error)
	UserByEmail(ctx context.Context, email string) (*models.User, error)
	SetRole(ctx context.Context, id string, role string) error
}

type Service struct {
	log        *slog.Logger
	repository Repository
}

func New(log *slog.Logger, repository Repository) *Service {
	return &Service{
		log:        log,
		repository: repository,
	}
}

func (s *Service) Register(ctx context.Context, user requests.Register) error {
	const op = "service.Register"
	log := s.log.With(slog.String("operation", op))

	log.Info("register user")
	log.Debug("user", slog.Any("user", user))

	if user.Password != user.RepeatPassword {
		log.Debug("passwords do not match", slog.String("password", user.Password), slog.String("repeat_password", user.RepeatPassword))
		return fmt.Errorf("%w: %s", ErrBadRequest, "passwords do not match")
	}

	log.Debug("creating user")

	hashedPassword, err := HashPassword(user.Password)
	if err != nil {
		log.Error("failed to hash password", err)
		return fmt.Errorf("%w: %s", ErrInternalServer, "failed to hash password")
	}

	log.Debug("saving user")

	err = s.repository.SaveUser(ctx, &models.User{
		ID:           uuid.NewString(),
		Email:        user.Email,
		Fullname:     user.Fullname,
		Password:     hashedPassword,
		Role:         user.Role,
		CurrentStage: user.CurrentStage,
	})

	if err != nil {
		log.Error("failed to save user", err)
		return fmt.Errorf("%w: %s", ErrInternalServer, "failed to save user")
	}

	return nil
}

func (s *Service) Login(ctx context.Context, user requests.Login) (string, error) {
	const op = "service.Login"
	log := s.log.With(slog.String("operation", op))

	log.Info("login user")
	log.Debug("user", slog.Any("user", user))

	log.Debug("getting user by email")

	userFind, err := s.repository.UserByEmail(ctx, user.Email)

	if err != nil {
		log.Error("failed to get user by email", err)
		return "", fmt.Errorf("%w: %s", ErrInternalServer, "failed to get user by email")
	}

	log.Debug("checking password")
	if !CheckPasswordHash(user.Password, userFind.Password) {
		log.Debug("passwords do not match", slog.String("password", user.Password), slog.String("user_password", userFind.Password))
		return "", fmt.Errorf("%w: %s", ErrBadRequest, "passwords do not match")
	}

	log.Debug("creating token")
	token := CreateToken(userFind)

	return token, nil
}

func (s *Service) SetRole(ctx context.Context, id string, role string) error {
	const op = "service.SetRole"
	log := s.log.With(slog.String("operation", op))

	log.Info("set role")

	if err := s.repository.SetRole(ctx, id, role); err != nil {
		log.Error("failed to set role", err)
		return fmt.Errorf("%w: %s", ErrInternalServer, "failed to set role")
	}

	return nil
}

func (s *Service) GetUser(ctx context.Context, id string) (*models.User, error) {
	const op = "service.GetUser"
	log := s.log.With(slog.String("operation", op))

	log.Info("get user")

	user, err := s.repository.User(ctx, id)

	if err != nil {
		log.Error("failed to get user", err)
		return nil, fmt.Errorf("%w: %s", ErrInternalServer, "failed to get user")
	}

	return user, nil
}

func (s *Service) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	const op = "service.GetUserByEmail"
	log := s.log.With(slog.String("operation", op))

	log.Info("get user by email")

	user, err := s.repository.UserByEmail(ctx, email)

	if err != nil {
		log.Error("failed to get user by email", err)
		return nil, fmt.Errorf("%w: %s", ErrInternalServer, "failed to get user by email")
	}

	return user, nil
}
