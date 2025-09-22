package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	RabbitURLServer string
	RabbitURLAgent  string
	DBHost          string
	DBPort          string
	DBUser          string
	DBPass          string
	DBName          string
}

func LoadConfig() *Config {
	// load .env file
	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file found: %v", err)
	}

	// default values
	cfg := &Config{
		RabbitURLServer: getEnv("RABBIT_URL_SERVER", "amqp://guest:guest@localhost:5672/"),
		RabbitURLAgent:  getEnv("RABBIT_URL_AGENT", "amqp://guest:guest@localhost:5672/"),
		DBHost:          getEnv("DB_HOST", "localhost"),
		DBPort:          getEnv("DB_PORT", "5432"),
		DBUser:          getEnv("DB_USER", "postgres"),
		DBPass:          getEnv("DB_PASS", "password"),
		DBName:          getEnv("DB_NAME", "monitoring"),
	}
	return cfg
}

func getEnv(key, defaultVal string) string {
	// check variable in the system environment
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}
