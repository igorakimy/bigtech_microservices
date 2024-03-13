package env

import (
	"github.com/igorakimy/bigtech_microservices/internal/config"
	"github.com/pkg/errors"
	"net"
	"os"
)

const (
	swaggerHostEnvName = "SWAGGER_HOST"
	swaggerPortEnvName = "SWAGGER_PORT"
)

type swaggerConfig struct {
	host string
	port string
}

func NewSwaggerConfig() (config.HTTPConfig, error) {
	swaggerHost := os.Getenv(swaggerHostEnvName)
	if len(swaggerHost) == 0 {
		return nil, errors.New("swagger host not found")
	}

	swaggerPort := os.Getenv(swaggerPortEnvName)
	if len(swaggerPort) == 0 {
		return nil, errors.New("swagger port not found")
	}

	return &swaggerConfig{
		host: swaggerHost,
		port: swaggerPort,
	}, nil
}

func (cfg *swaggerConfig) Address() string {
	return net.JoinHostPort(cfg.host, cfg.port)
}
