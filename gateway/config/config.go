package config

import (
	"os"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	ServerPort  string
	Environment string
	GinMode     string
	Domain      string
}

var Config *AppConfig

func LoadConfig() {

	if os.Getenv("ENVIRONMENT") != "production" {
		godotenv.Load()
	}

	Config = &AppConfig{
		ServerPort:  os.Getenv("PORT"),
		Environment: os.Getenv("ENVIRONMENT"),
		GinMode:     os.Getenv("GIN_MODE"),
		Domain:      os.Getenv("DOMAIN"),
	}

}
