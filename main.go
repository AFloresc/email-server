package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/resendlabs/resend-go"
)

type ContactRequest struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Message string `json:"message"`
}

func main() {
	http.HandleFunc("/contact", handleContact)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("Server running on port", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func handleContact(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req ContactRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if req.Name == "" || req.Email == "" || req.Message == "" {
		http.Error(w, "Missing fields", http.StatusBadRequest)
		return
	}

	if err := sendEmailResend(req); err != nil {
		log.Println("Resend error:", err)
		http.Error(w, "Failed to send email", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
	})
}

func sendEmailResend(req ContactRequest) error {
	apiKey := os.Getenv("RESEND_API_KEY")
	toEmail := os.Getenv("TO_EMAIL")

	client := resend.NewClient(apiKey)

	params := &resend.SendEmailRequest{
		From:    "Portfolio Contact <onboarding@resend.dev>",
		To:      []string{toEmail},
		Subject: "Nuevo mensaje de contacto",
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
