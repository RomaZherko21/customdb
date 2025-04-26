package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBPath      string
	Environment string
	Port        string
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	return &Config{
		DBPath:      os.Getenv("DB_PATH"),
		Environment: os.Getenv("ENV"),
		Port:        os.Getenv("PORT"),
	}, nil
}
