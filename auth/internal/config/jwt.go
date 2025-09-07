package config

import (
	"log"
	"os"
	"time"
)

type JWTConfig struct {
	RefreshSecretKey   string
	AccessSecretKey    string
	AccessTokenExpiry  time.Duration
	RefreshTokenExpiry time.Duration
}

func NewJWTConfig() *JWTConfig {
	refreshSecret := os.Getenv("JWT_REFRESH_SECRET")
	if refreshSecret == "" {
		log.Fatal("JWT_REFRESH_SECRET environment variable is required")
	}

	accessSecret := os.Getenv("JWT_ACCESS_SECRET")
	if accessSecret == "" {
		log.Fatal("JWT_ACCESS_SECRET environment variable is required")
	}

	accessTokenExpiry := getEnvDuration("JWT_ACCESS_TOKEN_EXPIRY")
	if accessTokenExpiry == 0 {
		log.Fatal("JWT_ACCESS_TOKEN_EXPIRY environment variable is required")
	}

	refreshTokenExpiry := getEnvDuration("JWT_REFRESH_TOKEN_EXPIRY")
	if refreshTokenExpiry == 0 {
		log.Fatal("JWT_REFRESH_TOKEN_EXPIRY environment variable is required")
	}

	return &JWTConfig{
		RefreshSecretKey:   refreshSecret,
		AccessSecretKey:    accessSecret,
		AccessTokenExpiry:  accessTokenExpiry,
		RefreshTokenExpiry: refreshTokenExpiry,
	}
}

func getEnvDuration(key string) time.Duration {
	if v := os.Getenv(key); v != "" {
		if duration, err := time.ParseDuration(v); err == nil {
			return duration
		}
	}
	return 0
}
