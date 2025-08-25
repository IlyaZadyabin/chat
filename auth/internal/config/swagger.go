package config

import (
	"net"
	"os"
)

const (
	swaggerHostEnvName = "SWAGGER_HOST"
	swaggerPortEnvName = "SWAGGER_PORT"
)

type SwaggerConfig interface {
	Address() string
}

type swaggerConfig struct {
	host string
	port string
}

func NewSwaggerConfig() SwaggerConfig {
	host := os.Getenv(swaggerHostEnvName)
	if len(host) == 0 {
		host = "localhost"
	}

	port := os.Getenv(swaggerPortEnvName)
	if len(port) == 0 {
		port = "8082"
	}

	return &swaggerConfig{
		host: host,
		port: port,
	}
}

func (cfg *swaggerConfig) Address() string {
	return net.JoinHostPort(cfg.host, cfg.port)
}
