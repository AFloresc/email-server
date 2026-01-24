package router

import (
	"fmt"
	"log/slog"
	"net/http"

	"email-server/internal/config"
	"email-server/internal/handlers"
	"email-server/internal/metrics"
	"email-server/internal/middleware"

	"github.com/rs/cors"
)

func StartServer(cfg *config.Config) error {
	mux := http.NewServeMux()

	mux.HandleFunc("/contact", handlers.ContactHandler(cfg))
	mux.HandleFunc("/metrics", metricsHandler)

	handler := middleware.Metrics(
		middleware.RateLimit(
			middleware.Cooldown(
				middleware.UserAgent(
					middleware.BodyLimit(mux),
				),
			),
		),
	)

	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   cfg.AllowedOrigins,
		AllowedMethods:   []string{"POST", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type"},
		AllowCredentials: false,
	}).Handler(handler)

	addr := ":" + cfg.Port
	slog.Info("Server listening", "addr", addr)

	return http.ListenAndServe(addr, corsHandler)
}

func metricsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")

	fmt.Fprintf(w, "total_requests %d\n", metrics.M.TotalRequests)
	fmt.Fprintf(w, "successful_requests %d\n", metrics.M.SuccessfulRequests)
	fmt.Fprintf(w, "validation_errors %d\n", metrics.M.ValidationErrors)
	fmt.Fprintf(w, "honeypot_blocks %d\n", metrics.M.HoneypotBlocks)
	fmt.Fprintf(w, "firewall_blocks %d\n", metrics.M.FirewallBlocks)
	fmt.Fprintf(w, "rate_limit_blocks %d\n", metrics.M.RateLimitBlocks)
	fmt.Fprintf(w, "cooldown_blocks %d\n", metrics.M.CooldownBlocks)
	fmt.Fprintf(w, "user_agent_blocks %d\n", metrics.M.UserAgentBlocks)
	fmt.Fprintf(w, "payload_too_large %d\n", metrics.M.PayloadTooLarge)
	fmt.Fprintf(w, "email_send_errors %d\n", metrics.M.EmailSendErrors)
	fmt.Fprintf(w, "total_processing_time_ms %d\n", metrics.M.TotalProcessingTime)
}
