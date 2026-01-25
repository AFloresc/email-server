package main

import (
	"log/slog"
	"os"

	"email-server/internal/config"
	emailstats "email-server/internal/emailstgats"
	"email-server/internal/logging"
	"email-server/internal/router"
)

func main() {
	logger := logging.NewLogger()
	slog.SetDefault(logger)

	cfg := config.Load()
	cfg.AppName = "Email Server"

	slog.Info("Starting contact server", "port", cfg.Port)

	// Inicializar el contador mensual
	if err := emailstats.Init("internal/data/email_stats.json"); err != nil {
		slog.Error("Failed to initialize email stats", "error", err)
		os.Exit(1)
	}

	if err := router.StartServer(cfg); err != nil {
		slog.Error("Server stopped with error", "error", err.Error())
		os.Exit(1)
	}
}
