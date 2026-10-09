package config

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
)
type Config struct {
	AppPort		string
	ClientOrigin string
	DatabaseURL string
	DBHost		string
	DBPort		string
	DBUser		string
	DBPassword	string
	DBName		string
	DBSSLMode	string
	JWTSecret string
	JWTExpiresIn time.Duration
	JWTIssuer string
	JWTAudience string
	MailtrapToken    string
  MailtrapFromEmail string
  MailtrapFromName  string
}

var Cfg *Config
func LoadConfig(){
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
	
	Cfg = &Config{
		AppPort:     port,
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

		JWTSecret:		 os.Getenv("JWT_SECRET"),
		JWTExpiresIn:	 expiresIn,
		JWTIssuer:		 os.Getenv("JWT_ISSUER"),
		JWTAudience:	 os.Getenv("JWT_AUDIENCE"),

		MailtrapToken:     os.Getenv("MAILTRAP_TOKEN"),
    MailtrapFromEmail: os.Getenv("MAILTRAP_FROM_EMAIL"),
    MailtrapFromName:  os.Getenv("MAILTRAP_FROM_NAME"),			
	}
}