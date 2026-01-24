package logging

import (
	"log/slog"
	"os"
)

// NewLogger crea un logger global en formato JSON.
// Este logger se usa en toda la aplicación mediante slog.SetDefault().
func NewLogger() *slog.Logger {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo, // puedes cambiar a Debug si quieres más detalle
	})

	return slog.New(handler)
}
