package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/resendlabs/resend-go"
	"github.com/rs/cors"
)

type ContactRequest struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Message string `json:"message"`
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/contact", handleContact)

	// CORS para desarrollo local
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowedMethods:   []string{"POST", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type"},
		AllowCredentials: true,
	}).Handler(mux)

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
