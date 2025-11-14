package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)


type Config struct {
	ServerAddress       string
	DBURL               string
	AppMode             string
	JWTSecret           string
	JWTExpiry           string
	RefreshTokenExpiry  string
	LogLevel            string
	LogFormat           string
	LogOutput           string
	MaxLoginAttempts    int
	LoginAttemptWindow  string
}

func New() *Config {
	// Set up logging to a file
	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Failed to open log file:", err)
	}
	log.SetOutput(file)
	if err := godotenv.Load(".env"); err != nil {
		log.Println("\nWarning: .env file not found, using system environment variables")
	}

	return &Config{
		ServerAddress:       getEnv("SERVER_ADDRESS", ":3333"),
		DBURL:               getEnv("DB_URL", "host=localhost user=doorman_user password=doorman_password dbname=doorman_db port=5432 sslmode=disable TimeZone=UTC"),
		AppMode:             getEnv("APP_MODE", "development"),
		JWTSecret:           getEnv("JWT_SECRET", "change-this-secret-key"),
		JWTExpiry:           getEnv("JWT_EXPIRY", "15m"),
		RefreshTokenExpiry:  getEnv("REFRESH_TOKEN_EXPIRY", "7d"),
		LogLevel:            getEnv("LOG_LEVEL", "info"),
		LogFormat:           getEnv("LOG_FORMAT", "json"),
		LogOutput:           getEnv("LOG_OUTPUT", "stdout"),
		MaxLoginAttempts:    5,
		LoginAttemptWindow:  getEnv("LOGIN_ATTEMPT_WINDOW", "15m"),
	}
}

func getEnv(key, defaultVal string) string {
	if val, exists := os.LookupEnv(key); exists {
		return val
	}
	return defaultVal
}
