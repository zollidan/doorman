package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerAddress string
	DBURL         string
	AppMode       string
}

func New() *Config {
	// Set up logging to a file
	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Failed to open log file:", err)
	}
	log.SetOutput(file)
	if err := godotenv.Load(".env"); err != nil {
		log.Println("\n⚠️ Warning: .env file not found, using system environment variables")
	}

	return &Config{
		ServerAddress: getEnv("SERVER_ADDRESS", "localhost:3333"),
		DBURL:         getEnv("DB_URL", "test.db"),
		AppMode:       getEnv("APP_MODE", "development"),
	}
}

func getEnv(key, defaultVal string) string {
	if val, exists := os.LookupEnv(key); exists {
		return val
	}
	return defaultVal
}
