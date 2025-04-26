package logger

import (
	"context"
	"path/filepath"
	"strconv"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func New() *zap.Logger {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}

	return logger
}

func Interceptor(logger *zap.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		next grpc.UnaryHandler,
	) (resp any, err error) {

		logger.Info(
			"new request", zap.String("method", info.FullMethod),
			zap.Any("request", req),
			zap.Time("time", time.Now()),
		)

		return next(ctx, req)
	}
}

func InterceptorLogger(l *zap.Logger) logging.Logger {
	return logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
		filds := make([]zap.Field, len(fields))
		for i, f := range fields {
			filds[i] = zap.Any(filepath.Base(strconv.Itoa(i)), f)
		}
		l.Info(msg, filds...)
	})
}
