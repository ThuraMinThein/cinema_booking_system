package config

import (
	"os"

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
	}

}
