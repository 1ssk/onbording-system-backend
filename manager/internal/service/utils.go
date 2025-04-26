package service

import (
	"context"
	"fmt"
	"time"

	test "github.com/glebkasilov/grpc-manager/pkg/api"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/metadata"
)

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func getCallerInfo(ctx context.Context) string {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		return fmt.Sprintf("user:%v", md.Get("user-agent"))
	}
	return "unknown"
}

func safeLogMeeting(m *test.Meeting) zap.Field {
	if m == nil {
		return zap.Skip()
	}
	return zap.Object("meeting", zapcore.ObjectMarshalerFunc(func(enc zapcore.ObjectEncoder) error {
		enc.AddString("id", m.Id)
		enc.AddString("user_id", m.UserId)
		enc.AddString("title", m.Title)
		return nil
	}))
}

func convertTime(t string) (time.Time, error) {
	tt, err := time.Parse("2006-01-02", t)
	if err != nil {
		return time.Time{}, err
	}
	return tt, nil
}
