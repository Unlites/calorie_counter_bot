package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	TelegramToken string
	DbUser        string
	DbPassword    string
}

func Init() (*Config, error) {
	var cfg Config
	if err := parseEnv(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func parseEnv(cfg *Config) error {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}

	TelegramToken, exists := os.LookupEnv("TOKEN")
	if !exists {
		log.Printf("No %s var found", TelegramToken)
	}
	DbUser, exists := os.LookupEnv("DBUSER")
	if !exists {
		log.Printf("No %s var found", DbUser)
	}
	DbPassword, exists := os.LookupEnv("DBPASSWORD")
	if !exists {
		log.Printf("No %s var found", DbPassword)
	}

	cfg.TelegramToken = TelegramToken
	cfg.DbUser = DbUser
	cfg.DbPassword = DbPassword

	return nil
}
