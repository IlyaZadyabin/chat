package config

import (
	"log"
	"os"
)

type AuthConfig struct {
	ServiceAddr string
}

func NewAuthConfig() *AuthConfig {
	serviceAddr := os.Getenv("AUTH_SERVICE_ADDR")
	if serviceAddr == "" {
		log.Fatal("AUTH_SERVICE_ADDR environment variable is required")
	}

	return &AuthConfig{
		ServiceAddr: serviceAddr,
	}
}
