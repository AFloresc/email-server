package config

import (
	"os"
	"strings"
)

type Config struct {
	Port           string
	AllowedOrigins []string
	ResendAPIKey   string
	ToEmail        string
	AppName        string
}

func Load() *Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	originsEnv := os.Getenv("ALLOWED_ORIGINS")
	var origins []string
	if originsEnv != "" {
		origins = strings.Split(originsEnv, ",")
	} else {
		origins = []string{"http://localhost:5173"}
	}

	return &Config{
		Port:           port,
		AllowedOrigins: origins,
		ResendAPIKey:   os.Getenv("RESEND_API_KEY"),
		ToEmail:        os.Getenv("TO_EMAIL"),
	}
}
