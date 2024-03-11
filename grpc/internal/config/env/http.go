package env

import (
	"github.com/igorakimy/bigtech_microservices/internal/config"
	"github.com/pkg/errors"
	"net"
	"os"
)

const (
	httpHostEnvName = "HTTP_HOST"
	httpPortEnvName = "HTTP_PORT"
)

type httpConfig struct {
	host string
	port string
}

func NewHTTPConfig() (config.HTTPConfig, error) {
	httpHost := os.Getenv(httpHostEnvName)
	if len(httpHost) == 0 {
		return nil, errors.New("http host not found")
	}

	httpPort := os.Getenv(httpPortEnvName)
	if len(httpPort) == 0 {
		return nil, errors.New("http port not found")
	}

	return &httpConfig{
		host: httpHost,
		port: httpPort,
	}, nil
}

func (cfg *httpConfig) Address() string {
	return net.JoinHostPort(cfg.host, cfg.port)
}
