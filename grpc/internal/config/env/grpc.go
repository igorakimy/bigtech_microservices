package env

import (
	"errors"
	"github.com/igorakimy/bigtech_microservices/internal/config"
	"net"
	"os"
)

const (
	grpcHostEnvName = "GRPC_HOST"
	grpcPortEnvName = "GRPC_PORT"
)

type grpcConfig struct {
	host string
	port string
}

func NewGRPCConfig() (config.GRPCConfig, error) {
	grpcHost := os.Getenv(grpcHostEnvName)
	if len(grpcHost) == 0 {
		return nil, errors.New("grpc host not found")
	}

	grpcPort := os.Getenv(grpcPortEnvName)
	if len(grpcPort) == 0 {
		return nil, errors.New("grpc port not found")
	}

	return &grpcConfig{
		host: grpcHost,
		port: grpcPort,
	}, nil
}

func (cfg *grpcConfig) Address() string {
	return net.JoinHostPort(cfg.host, cfg.port)
}
