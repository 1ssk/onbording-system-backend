package application

import (
	"log/slog"

	"github.com/glebkasilov/authorization/internal/server"
	"github.com/glebkasilov/authorization/internal/service"
	"github.com/glebkasilov/authorization/internal/storage"
)

type Aplication struct {
	server *server.Server
}

func New(log *slog.Logger) *Aplication {
	storage := storage.New()
	service := service.New(log, storage)
	server := server.New(service)
	return &Aplication{
		server: server,
	}
}

func (a *Aplication) Start() {
	a.server.Start()
}

func (a *Aplication) GracefulStop() {
	a.server.GracefulStop()
}
