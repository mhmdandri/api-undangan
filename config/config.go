package config

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort      string
	ClientOrigin string
	DatabaseURL  string
	DBHost       string
	DBPort       string
	DBUser       string
	DBPassword   string
	DBName       string
	DBSSLMode    string
	JWTSecret    string
	JWTExpiresIn time.Duration
	JWTIssuer    string
	JWTAudience  string

	// Legacy / Fallback Mailtrap
	MailtrapToken     string
	MailtrapFromEmail string
	MailtrapFromName  string

	// New Email Service (SMTP iCloud API)
	EmailAPIURL string
	EmailAPIKey string
}

var Cfg *Config

func LoadConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, use environment variables instead")
	}
	expiresIn, err := time.ParseDuration(os.Getenv("JWT_EXPIRES_IN"))
	if err != nil {
		expiresIn = time.Hour * 24
	}

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = ":8888"
	} else if !strings.HasPrefix(port, ":") {
		port = ":" + port
	}

	emailAPIURL := os.Getenv("EMAIL_API_URL")
	if emailAPIURL == "" {
		emailAPIURL = "http://127.0.0.1:8082/api/send-email"
	}

	Cfg = &Config{
		AppPort: port,
		ClientOrigin: func() string {
			if v := os.Getenv("CLIENT_ORIGIN"); v != "" {
				return v
			}
			return "https://weddingofandricica.me"
		}(),
		DatabaseURL: os.Getenv("DATABASE_URL"),
		DBHost:      os.Getenv("DB_HOST"),
		DBPort:      os.Getenv("DB_PORT"),
		DBUser:      os.Getenv("DB_USER"),
		DBPassword:  os.Getenv("DB_PASSWORD"),
		DBName:      os.Getenv("DB_NAME"),
		DBSSLMode:   os.Getenv("DB_SSLMODE"),

		JWTSecret:    os.Getenv("JWT_SECRET"),
		JWTExpiresIn: expiresIn,
		JWTIssuer:    os.Getenv("JWT_ISSUER"),
		JWTAudience:  os.Getenv("JWT_AUDIENCE"),

		MailtrapToken:     os.Getenv("MAILTRAP_TOKEN"),
		MailtrapFromEmail: os.Getenv("MAILTRAP_FROM_EMAIL"),
		MailtrapFromName:  os.Getenv("MAILTRAP_FROM_NAME"),

		EmailAPIURL: emailAPIURL,
		EmailAPIKey: os.Getenv("EMAIL_API_KEY"),
	}
}