package main

import (
	"log/slog"
	"os"

	"email-server/internal/config"
	"email-server/internal/logging"
	"email-server/internal/router"
)

func main() {
	logger := logging.NewLogger()
	slog.SetDefault(logger)

	cfg := config.Load()

	slog.Info("Starting contact server", "port", cfg.Port)

	if err := router.StartServer(cfg); err != nil {
		slog.Error("Server stopped with error", "error", err.Error())
		os.Exit(1)
	}
}
