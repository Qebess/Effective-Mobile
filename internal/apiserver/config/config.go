package config

import (
	"fmt"
	"log"
	"os"
)

type Config struct {
	AppEnv  string
	AppPort string
	ConnStr string
}

func MustLoad() *Config {
	appEnv := os.Getenv("LOG_LEVEL")

	if appEnv == "" {
		log.Fatal("LOG_LEVEL is not set")
	}

	appPort := os.Getenv("APP_PORT")
	if appPort == "" {
		log.Fatal("APP_PORT is not set")
	}

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	cfg := &Config{AppEnv: appEnv, AppPort: appPort, ConnStr: connStr}

	return cfg
}
