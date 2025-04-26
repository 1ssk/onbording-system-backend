package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/zap"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"notification/internal/api"
	"notification/internal/config"
	"notification/internal/logger"
	"notification/internal/service"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	l := logger.New()

	server_cfg, err := config.New()
	if err != nil {
		l.Fatal("failed to read config", zap.Error(err))
	}

	srv := service.New(server_cfg, l)

	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", server_cfg.Host, server_cfg.NotificationsGrpcPort))
	if err != nil {
		l.Fatal("failed to listen", zap.Error(err))
	}

	recoveryOpts := []recovery.Option{recovery.WithRecoveryHandler(func(p interface{}) (err error) {
		return err
	})}

	go runRest(ctx, l, server_cfg)
	server := grpc.NewServer(grpc.ChainUnaryInterceptor(recovery.UnaryServerInterceptor(recoveryOpts...), logging.UnaryServerInterceptor(logger.InterceptorLogger(l))))

	ctx, stop := signal.NotifyContext(ctx, os.Interrupt)
	defer stop()

	l.Info("Server started", zap.String("port", server_cfg.Port))

	go func() {
		select {
		case <-ctx.Done():
			server.GracefulStop()
			l.Info("Server stopped")
			<-ctx.Done()
		}
	}()

	api.RegisterNotificationServiceServer(server, srv)

	if err := server.Serve(lis); err != nil {
		l.Fatal("failed to serve", zap.Error(err))
	}
}

func runRest(ctx context.Context, logger *zap.Logger, server_cfg *config.Config) {
	rt := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	grpcEndpoint := fmt.Sprintf("%s:%d", server_cfg.Host, server_cfg.NotificationsGrpcPort)

	err := api.RegisterNotificationServiceHandlerFromEndpoint(ctx, rt, grpcEndpoint, opts)
	if err != nil {
		logger.Error("Failed to register gateway", zap.Error(err))
		return
	}

	h2Server := &http2.Server{}
	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", server_cfg.Port),
		Handler: h2c.NewHandler(rt, h2Server),
	}

	logger.Info("REST server started")
	if err := server.ListenAndServe(); err != nil {
		logger.Error("Failed to start server", zap.Error(err))
	}
}
