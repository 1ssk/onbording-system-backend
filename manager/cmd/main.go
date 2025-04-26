package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/glebkasilov/grpc-manager/internal/config"
	"github.com/glebkasilov/grpc-manager/internal/service"

	test "github.com/glebkasilov/grpc-manager/pkg/api"
	"github.com/glebkasilov/grpc-manager/pkg/database/postgres"
	"github.com/glebkasilov/grpc-manager/pkg/logger"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	ctx := context.Background()
	mainLogger := logger.New()

	cfg := config.Config()

	server_cfg := cfg.Server
	datebase_cfg := cfg.Database

	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", server_cfg.Host, server_cfg.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	storage, err := postgres.New(&datebase_cfg)
	if err != nil {
		log.Fatalf("failed to create storage: %v", err)
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	if err := redisClient.Ping(ctx).Err(); err != nil {
		mainLogger.Error("Failed to connect to Redis", zap.Error(err))
	}

	srv := service.New(*storage, redisClient, mainLogger)

	recoveryOpts := []recovery.Option{recovery.WithRecoveryHandler(func(p interface{}) (err error) {
		return err
	})}
	go runRest(ctx, mainLogger, &server_cfg)

	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			recovery.UnaryServerInterceptor(recoveryOpts...),
			logging.UnaryServerInterceptor(logger.InterceptorLogger(mainLogger)),
			CacheInterceptor(redisClient, cfg.Redis.TTL),
		),
	)

	ctx, stop := signal.NotifyContext(ctx, os.Interrupt)
	defer stop()

	mainLogger.Info("Starting server", zap.String("server", fmt.Sprintf("%s:%d", server_cfg.Host, server_cfg.Port)))

	server = grpc.NewServer(grpc.ChainUnaryInterceptor(recovery.UnaryServerInterceptor(recoveryOpts...), logging.UnaryServerInterceptor(logger.InterceptorLogger(mainLogger))))

	ctx, stop = signal.NotifyContext(ctx, os.Interrupt)
	defer stop()

	go func() {
		select {
		case <-ctx.Done():
			server.GracefulStop()
			mainLogger.Info("Server stopped")
			<-ctx.Done()
		}
	}()

	test.RegisterMeetingServiceServer(server, srv)

	if err := server.Serve(lis); err != nil {
		mainLogger.Error("Failed to start server", zap.Error(err))
	}
}

func runRest(ctx context.Context, logger *logger.Logger, server_cfg *config.ServerCfg) {
	rt := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	grpcEndpoint := fmt.Sprintf("%s:%d", server_cfg.Host, server_cfg.Port)

	err := test.RegisterMeetingServiceHandlerFromEndpoint(ctx, rt, grpcEndpoint, opts)
	if err != nil {
		logger.Error("Failed to register gateway", zap.Error(err))
		return
	}

	h2Server := &http2.Server{}
	server := &http.Server{
		Addr:    fmt.Sprintf(":%v", server_cfg.HttpPort),
		Handler: h2c.NewHandler(rt, h2Server),
	}

	logger.Info("Starting REST server", zap.Int("port", server_cfg.HttpPort))
	if err := server.ListenAndServe(); err != nil {
		logger.Error("Failed to start REST server", zap.Error(err))
	}
}

func CacheInterceptor(rdb *redis.Client, ttl time.Duration) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		cacheKey := info.FullMethod + ":" + hashRequest(req)

		var response interface{}
		if err := getFromCache(ctx, rdb, cacheKey, &response); err == nil {
			return response, nil
		}

		res, err := handler(ctx, req)
		if err != nil {
			return nil, err
		}

		setToCache(ctx, rdb, cacheKey, res, ttl)

		return res, nil
	}
}

func hashRequest(req interface{}) string {
	data, _ := json.Marshal(req)
	return string(data)
}

func getFromCache(ctx context.Context, rdb *redis.Client, key string, dest interface{}) error {
	data, err := rdb.Get(ctx, key).Bytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(data, dest)
}

func setToCache(ctx context.Context, rdb *redis.Client, key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return rdb.Set(ctx, key, data, ttl).Err()
}
