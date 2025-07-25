package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func NewConfig() *Config {
	return &Config{
		Host:     getEnv("PG_HOST", "localhost"),
		Port:     getEnv("PG_PORT", "5432"),
		User:     getEnv("AUTH_POSTGRES_USER", "auth_user"),
		Password: getEnv("AUTH_POSTGRES_PASSWORD", "auth_password"),
		DBName:   getEnv("AUTH_POSTGRES_DB", "auth"),
		SSLMode:  getEnv("PG_SSLMODE", "disable"),
	}
}

func (c *Config) GetDSN() string {
	return fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=%s",
		c.Host, c.Port, c.DBName, c.User, c.Password, c.SSLMode)
}

func NewConnection(cfg *Config) (*pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.Connect(ctx, cfg.GetDSN())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err = pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Successfully connected to database (pgx)")
	return pool, nil
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
