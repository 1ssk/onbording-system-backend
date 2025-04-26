package main

import (
	"log/slog"
	"os"
	"os/signal"

	"github.com/glebkasilov/authorization/internal/application"
	"github.com/glebkasilov/authorization/internal/config"
)

func main() {
	config.LoadConfig()
	h := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})

	mainLogger := slog.New(h)

	app := application.New(mainLogger)

	app.Start()

	ch := make(chan os.Signal, 1)

	signal.Notify(ch, os.Interrupt)

	<-ch

	app.GracefulStop()
}
