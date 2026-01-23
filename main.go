package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"

	"github.com/resendlabs/resend-go"
	"github.com/rs/cors"
	"golang.org/x/time/rate"
)

type ContactRequest struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Message string `json:"message"`
	Website string `json:"website"` // honeypot
}

// Rate limit variables
var visitors = make(map[string]*rate.Limiter)
var mu sync.Mutex

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/contact", handleContact)

	handler := rateLimitMiddleware(mux)
	allowedOrigins := strings.Split(os.Getenv("ALLOWED_ORIGINS"), ",")

	// CORS para desarrollo local y deploy
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowedMethods:   []string{"POST", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type"},
		AllowCredentials: false,
	}).Handler(handler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("Server running on port", port)
	log.Fatal(http.ListenAndServe(":"+port, corsHandler))
}

func handleContact(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req ContactRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := validateContact(req); err != nil {
		log.Println("Validation error:", err)
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if req.Website != "" {
		// Honeypot actrivated → bot detected
		log.Println("Honeypot triggered — bot blocked")
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	if err := sendEmailResend(req); err != nil {
		log.Println("Failed to send email:", err)
		http.Error(w, "Failed to send email", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status":"ok"}`))
}

func sendEmailResend(req ContactRequest) error {
	apiKey := os.Getenv("RESEND_API_KEY")
	toEmail := os.Getenv("TO_EMAIL")

	client := resend.NewClient(apiKey)

	htmlBody := fmt.Sprintf(`
		<div style="font-family: 'Segoe UI', Roboto, sans-serif; max-width: 600px; margin: auto; padding: 20px; background-color: #f9f9f9; border-radius: 12px; border: 1px solid #e0e0e0;">
			<h2 style="color: #0070f3; margin-bottom: 16px;">Nuevo mensaje desde tu portfolio</h2>

			<p style="margin: 8px 0;"><strong>Nombre:</strong> %s</p>
			<p style="margin: 8px 0;"><strong>Email:</strong> %s</p>

			<p style="margin: 8px 0;"><strong>Mensaje:</strong></p>
			<div style="margin: 12px 0; padding: 12px; background-color: #ffffff; border-left: 4px solid #0070f3;">
				<p style="margin: 0; white-space: pre-line;">%s</p>
			</div>

			<hr style="margin: 24px 0; border: none; border-top: 1px solid #ddd;" />

			<p style="font-size: 12px; color: #888;">
				Este mensaje fue enviado desde el formulario de contacto de tu portfolio.
			</p>
		</div>
	`, req.Name, req.Email, req.Message)

	params := &resend.SendEmailRequest{
		From:    "Portfolio Contact <onboarding@resend.dev>",
		To:      []string{toEmail},
		ReplyTo: req.Email,
		Subject: "Nuevo mensaje de contacto",
		Html:    htmlBody,
		Text: fmt.Sprintf(
			"Nombre: %s\nEmail: %s\nMensaje:\n%s",
			req.Name, req.Email, req.Message,
		),
	}

	res, err := client.Emails.Send(params)
	if err != nil {
		log.Printf("Resend error: %v\n", err)
		return err
	}

	log.Printf("Resend response: %+v\n", res)
	return nil
}

// Rate limit function
func getVisitor(ip string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	limiter, exists := visitors[ip]
	if !exists {
		limiter = rate.NewLimiter(1, 3) // 1 request/sec, burst of 3
		visitors[ip] = limiter
	}
	return limiter
}

func rateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr

		limiter := getVisitor(ip)
		if !limiter.Allow() {
			http.Error(w, "Too many requests", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// Fields validator
func validateContact(req ContactRequest) error {
	// Honeypot
	if req.Website != "" {
		return errors.New("bot detected")
	}

	// Name
	req.Name = strings.TrimSpace(req.Name)
	if len(req.Name) < 2 || len(req.Name) > 80 {
		return errors.New("invalid name")
	}

	// Email
	req.Email = strings.TrimSpace(req.Email)
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(req.Email) {
		return errors.New("invalid email")
	}

	// Message
	req.Message = strings.TrimSpace(req.Message)
	if len(req.Message) < 10 || len(req.Message) > 2000 {
		return errors.New("invalid message")
	}

	return nil
}
