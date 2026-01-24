package email

import (
	"fmt"
	"log/slog"

	"email-server/internal/config"

	"github.com/resendlabs/resend-go"
)

func SendEmailResend(cfg *config.Config, name, emailAddr, message string) error {
	client := resend.NewClient(cfg.ResendAPIKey)

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
	`, name, emailAddr, message)

	params := &resend.SendEmailRequest{
		From:    "Portfolio Contact <onboarding@resend.dev>",
		To:      []string{cfg.ToEmail},
		ReplyTo: emailAddr,
		Subject: "Nuevo mensaje de contacto",
		Html:    htmlBody,
		Text: fmt.Sprintf(
			"Nombre: %s\nEmail: %s\nMensaje:\n%s",
			name, emailAddr, message,
		),
	}

	res, err := client.Emails.Send(params)
	if err != nil {
		slog.Error("Resend error", "error", err.Error())
		return err
	}

	slog.Info("Resend response", "response", res)
	return nil
}
