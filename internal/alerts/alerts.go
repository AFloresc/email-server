package alerts

import (
	"fmt"
	"net/smtp"
)

func SendAlert(subject, body string) error {
	if cfg.SMTPHost == "" {
		return fmt.Errorf("alerts not initialized: missing SMTP config")
	}

	auth := smtp.PlainAuth("", cfg.Username, cfg.Password, cfg.SMTPHost)

	msg := "From: " + cfg.From + "\n" +
		"To: " + cfg.To + "\n" +
		"Subject: " + subject + "\n\n" +
		body

	addr := cfg.SMTPHost + ":" + cfg.SMTPPort

	return smtp.SendMail(addr, auth, cfg.From, []string{cfg.To}, []byte(msg))
}

func AlertMassSending(ip string, count int) error {
	subject := "ALERTA: Posible envío masivo sospechoso"
	body := fmt.Sprintf(
		"La IP %s ha enviado %d correos en poco tiempo. Esto podría indicar abuso.",
		ip, count,
	)
	return SendAlert(subject, body)
}

func AlertTierUsage(current int, limit int) error {
	subject := "ALERTA: Consumo de correos cercano al límite mensual"
	body := fmt.Sprintf(
		"Has enviado %d de %d correos este mes. Estás cerca del límite mensual.",
		current, limit,
	)
	return SendAlert(subject, body)
}

func AlertSMTPAnomalies(errCount int) error {
	subject := "ALERTA: Errores SMTP anómalos"
	body := fmt.Sprintf(
		"Se han detectado %d errores SMTP en los últimos minutos. Revisa el servidor.",
		errCount,
	)
	return SendAlert(subject, body)
}

func AlertRelayAttempt(ip string) error {
	subject := "ALERTA: Intento de relay no autorizado"
	body := fmt.Sprintf(
		"La IP %s ha intentado usar el servidor como relay sin permisos.",
		ip,
	)
	return SendAlert(subject, body)
}

func AlertMaliciousIP(ip string, reason string) error {
	subject := "ALERTA: IP sospechosa detectada"
	body := fmt.Sprintf(
		"La IP %s ha sido marcada como sospechosa. Motivo: %s",
		ip, reason,
	)
	return SendAlert(subject, body)
}
