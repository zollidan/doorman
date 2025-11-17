package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerAddress      string
	POSTGRES_USER      string
	POSTGRES_PASSWORD  string
	POSTGRES_DB        string
	POSTGRES_HOST      string
	POSTGRES_PORT      string
	AppMode            string
	JWTSecret          string
	JWTExpiry          string
	RefreshTokenExpiry string
	LogLevel           string
	LogFormat          string
	LogOutput          string
	MaxLoginAttempts   int
}

func New() *Config {

	if err := godotenv.Load(".env"); err != nil {
		log.Println("\nWarning: .env file not found, using system environment variables")
	}

	return &Config{
		ServerAddress:      getEnv("SERVER_ADDRESS", ":2222"),
		AppMode:            getEnv("APP_MODE", "development"),
		POSTGRES_USER:      getEnv("POSTGRES_USER", "doorman_user"),
		POSTGRES_PASSWORD:  getEnv("POSTGRES_PASSWORD", "doorman_password"),
		POSTGRES_DB:        getEnv("POSTGRES_DB", "doorman_db"),
		POSTGRES_HOST:      getEnv("POSTGRES_HOST", "localhost"),
		POSTGRES_PORT:      getEnv("POSTGRES_PORT", "5432"),
		JWTSecret:          getEnv("JWT_SECRET", "change-this-secret-key"),
		JWTExpiry:          getEnv("JWT_EXPIRY", "15m"),
		RefreshTokenExpiry: getEnv("REFRESH_TOKEN_EXPIRY", "7d"),
		LogLevel:           getEnv("LOG_LEVEL", "info"),
		LogFormat:          getEnv("LOG_FORMAT", "json"),
		LogOutput:          getEnv("LOG_OUTPUT", "stdout"),
		MaxLoginAttempts:   5,
	}
}

func getEnv(key, defaultVal string) string {
	if val, exists := os.LookupEnv(key); exists {
		return val
	}
	return defaultVal
}
