package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/glebkasilov/authorization/internal/config"
	"github.com/glebkasilov/authorization/internal/domain/models"
	"github.com/glebkasilov/authorization/internal/domain/requests"
	"github.com/glebkasilov/authorization/internal/server/routers"
)

type Service interface {
	Login(ctx context.Context, user requests.Login) (string, error)
	Register(ctx context.Context, user requests.Register) error
	SetRole(ctx context.Context, id string, role string) error
	GetUser(ctx context.Context, id string) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
}

type Server struct {
	server *http.Server
	engine *gin.Engine
}

func New(
	service Service,
) *Server {
	cfg := config.Config().Server

	r := gin.Default()

	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Port),
		Handler: r,
	}

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	group := r.Group("/api")
	routers.Register(group, service)

	return &Server{
		server: httpServer,
		engine: r,
	}
}

func (s *Server) Start() {
	if err := s.server.ListenAndServe(); err != nil {
		panic(err)
	}
}

func (s *Server) GracefulStop() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.server.Shutdown(ctx); err != nil {
		panic(err)
	}
}
