package main

import (
	"fmt"
	"github.com/joho/godotenv"
)

type Config struct{}

func readConfigFromEnv() (*Config, error) {
	cfg := new(Config)

	return cfg, nil
}

func readConfigFromDotEnv(path string) (*Config, error) {
	if err := godotenv.Load(path); err != nil {
		return nil, fmt.Errorf("failed to load .env file: %w", err)
	}

	return readConfigFromEnv()
}
