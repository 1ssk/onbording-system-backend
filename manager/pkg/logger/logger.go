package logger

import (
	"context"
	"os"
	"path/filepath"
	"strconv"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	l *zap.Logger
}

func New() *Logger {
	ws := zapcore.AddSync(os.Stdout)
	logger := zap.New(zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		ws,
		zap.DebugLevel,
	))
	return &Logger{
		l: logger,
	}
}

func (l *Logger) Info(msg string, fields ...zap.Field) {
	fields = append(fields, zap.String("level", "info"))
	l.l.Info(msg, fields...)
}

func (l *Logger) Error(msg string, fields ...zap.Field) {
	fields = append(fields, zap.String("level", "error"))
	l.l.Error(msg, fields...)
}

func (l *Logger) Fatal(msg string, fields ...zap.Field) {
	fields = append(fields, zap.String("level", "fatal"))
	l.l.Fatal(msg, fields...)
}

func (l *Logger) Debug(msg string, fields ...zap.Field) {
	fields = append(fields, zap.String("level", "debug"))
	l.l.Debug(msg, fields...)
}

func InterceptorLogger(l *Logger) logging.Logger {
	return logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
		filds := make([]zap.Field, len(fields))
		for i, f := range fields {
			filds[i] = zap.Any(filepath.Base(strconv.Itoa(i)), f)
		}
		l.l.Info(msg, filds...)
	})
}
