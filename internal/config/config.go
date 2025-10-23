package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	RabbitURL string
	GRPCPort  string
	DBHost    string
	DBPort    string
	DBUser    string
	DBPass    string
	DBName    string
}

// load configuration from .env file and environment variables
func LoadConfig() *Config {
	// load .env file
	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file found: %v", err)
	}

	// default values
	cfg := &Config{
		RabbitURL: getEnv("RABBIT_URL", "amqp://metric:metric@localhost:5672/"),
		GRPCPort:  getEnv("GRPC_PORT", "50051"),
		DBHost:    getEnv("DB_HOST", "localhost"),
		DBPort:    getEnv("DB_PORT", "5432"),
		DBUser:    getEnv("DB_USER", "postgres"),
		DBPass:    getEnv("DB_PASS", "postgres"),
		DBName:    getEnv("DB_NAME", "monitoring"),
	}
	return cfg
}

// check variable in the system environment
func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}
