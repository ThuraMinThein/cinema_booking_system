package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	ServerPort             string
	UsersServiceAddress    string
	BookingsServiceAddress string
	SeatsServiceAddress    string
	Environment            string
	GinMode                string
	Domain                 string

	RedisURL      string
	RedisUsername string
	RedisPassword string
	RedisDB       int
}

var Config *AppConfig

func LoadConfig() {

	if os.Getenv("ENVIRONMENT") != "production" {
		godotenv.Load()
	}

	Config = &AppConfig{
		ServerPort:             getEnv("PORT", "8080"),
		Environment:            getEnv("ENVIRONMENT", "development"),
		GinMode:                getEnv("GIN_MODE", "debug"),
		Domain:                 getEnv("DOMAIN", "localhost"),
		UsersServiceAddress:    getEnv("USERS_SERVICE_ADDRESS", "localhost:8081"),
		BookingsServiceAddress: getEnv("BOOKINGS_SERVICE_ADDRESS", "localhost:8082"),
		SeatsServiceAddress:    getEnv("SEATS_SERVICE_ADDRESS", "localhost:8083"),

		RedisURL:      getEnv("REDIS_URL", ""),
		RedisUsername: getEnv("REDIS_USERNAME", ""),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		RedisDB:       getEnvAsInt("REDIS_DB", 0),
	}

}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}
