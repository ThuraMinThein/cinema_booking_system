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

	// Database
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string

	// Redis
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
		ServerPort:             os.Getenv("PORT"),
		Environment:            os.Getenv("ENVIRONMENT"),
		GinMode:                os.Getenv("GIN_MODE"),
		Domain:                 os.Getenv("DOMAIN"),
		UsersServiceAddress:    os.Getenv("USERS_SERVICE_ADDRESS"),
		BookingsServiceAddress: os.Getenv("BOOKINGS_SERVICE_ADDRESS"),
		SeatsServiceAddress:    os.Getenv("SEATS_SERVICE_ADDRESS"),

		DBHost:     os.Getenv("DATABASE_HOST"),
		DBPort:     os.Getenv("DATABASE_PORT"),
		DBUser:     os.Getenv("DATABASE_USERNAME"),
		DBPassword: os.Getenv("DATABASE_PASSWORD"),
		DBName:     os.Getenv("DATABASE_NAME"),

		RedisURL:      os.Getenv("REDIS_URL"),
		RedisUsername: os.Getenv("REDIS_USERNAME"),
		RedisPassword: os.Getenv("REDIS_PASSWORD"),
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
