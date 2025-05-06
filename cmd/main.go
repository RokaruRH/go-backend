package main

import (
	"fmt"
	"log/slog"
	"os"
	"url-shortener/internal/config"
	"url-shortener/internal/lib/logger/sl"
	"url-shortener/internal/storage/sqlite"
)

const (
	envLocal      = "local"
	envDev        = "dev"
	envProduction = "production"
)

func main() {
	config := config.MustLoad()

	log := setuplogger(config.Env)

	log.Info("My server is started!!!!")

	storage, err := sqlite.New(config.StoragePath)
	if err != nil {
		log.Error("failed to init storage", sl.Err(err))
		os.Exit(1)
	}

	fmt.Println(storage.DeleteURL("g"))

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
