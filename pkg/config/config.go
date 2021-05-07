package config

import (
	"errors"
	"os"

	_ "github.com/joho/godotenv/autoload"
)

// GoblinatorConfig from .env
type GoblinatorConfig struct {
	BlizzardClientId     string
	BlizzardClientSecret string
}

func New() (*GoblinatorConfig, error) {
	config := GoblinatorConfig{}
	err := config.read()
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func (config *GoblinatorConfig) read() error {
	clientId := os.Getenv("BLIZZARD_CLIENT_ID")
	if clientId == "" {
		return errors.New("no client id in .env")
	}
	clientSecret := os.Getenv("BLIZZARD_CLIENT_SECRET")
	if clientSecret == "" {
		return errors.New("no client secret in .env")
	}
	config.BlizzardClientId = clientId
	config.BlizzardClientSecret = clientSecret
	return nil
}
