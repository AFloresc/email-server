package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"email-server/internal/config"
	"email-server/internal/email"
	"email-server/internal/metrics"
	"email-server/internal/security"
)

type ContactRequest struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Message string `json:"message"`
	Website string `json:"website"` // honeypot
}

func (c ContactRequest) GetName() string    { return c.Name }
func (c ContactRequest) GetEmail() string   { return c.Email }
func (c ContactRequest) GetMessage() string { return c.Message }
func (c ContactRequest) GetWebsite() string { return c.Website }

func ContactHandler(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		if r.Method != http.MethodPost {
			slog.Warn("Invalid method", "method", r.Method)
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		ipHash := security.HashIP(r.RemoteAddr)

		var req ContactRequest
		dec := json.NewDecoder(r.Body)
		if err := dec.Decode(&req); err != nil {
			if err.Error() == "http: request body too large" {
				slog.Warn("Payload too large", "ipHash", ipHash)
				metrics.Inc(&metrics.M.PayloadTooLarge)
				http.Error(w, "Payload too large", http.StatusRequestEntityTooLarge)
				return
			}

			slog.Warn("Invalid JSON body", "error", err.Error(), "ipHash", ipHash)
			metrics.Inc(&metrics.M.ValidationErrors)
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if err := security.ValidateContact(req); err != nil {
			if err.Error() == "bot detected" {
				slog.Warn("Honeypot triggered", "ipHash", ipHash)
				metrics.Inc(&metrics.M.HoneypotBlocks)
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}
			if err.Error() == "malicious pattern detected" {
				slog.Warn("Attack pattern blocked", "ipHash", ipHash)
				metrics.Inc(&metrics.M.FirewallBlocks)
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}

			slog.Warn("Validation error", "error", err.Error(), "ipHash", ipHash)
			metrics.Inc(&metrics.M.ValidationErrors)
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}

		if err := email.SendEmailResend(cfg, req.Name, req.Email, req.Message); err != nil {
			slog.Error("Failed to send email", "error", err.Error(), "ipHash", ipHash)
			metrics.Inc(&metrics.M.EmailSendErrors)
			http.Error(w, "Failed to send email", http.StatusInternalServerError)
			return
		}

		slog.Info("Email sent successfully", "ipHash", ipHash, "email", req.Email)
		metrics.Inc(&metrics.M.SuccessfulRequests)

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))

		_ = time.Now() // placeholder si quieres añadir más métricas por request
	}
}
