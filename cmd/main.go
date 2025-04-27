package main

import (
	"log/slog"
	"os"
	"url-shortener/internal/config"
)

const (
	envLocal      = "local"
	envDev        = "dev"
	envProduction = "production"
)

func main() {
	config := config.MustLoad()

	logger := setuplogger(config.Env)

	logger.Info("My server is started!!!!")

	//TODO init storage: sqlite

	//TODO init router: chi "chi reader"

	//TODO init run server
}

func setuplogger(environment string) *slog.Logger {
	var log *slog.Logger
	switch environment {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProduction:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}
	return log
}
