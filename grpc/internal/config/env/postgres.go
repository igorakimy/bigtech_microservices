package env

import (
	"errors"
	"fmt"
	"github.com/igorakimy/bigtech_microservices/internal/config"
	"os"
)

const (
	hostEnvName     = "POSTGRES_HOST"
	portEnvName     = "POSTGRES_PORT"
	dbEnvName       = "POSTGRES_DBNAME"
	usernameEnvName = "POSTGRES_USER"
	passwordEnvName = "POSTGRES_PASSWORD"
)

type postgresConfig struct {
	host     string
	port     string
	dbname   string
	username string
	password string
}

func NewPostgresConfig() (config.PostgresConfig, error) {
	host := os.Getenv(hostEnvName)
	if len(host) == 0 {
		return nil, errors.New("postgres host not found")
	}

	port := os.Getenv(portEnvName)
	if len(port) == 0 {
		return nil, errors.New("postgres port not found")
	}

	dbName := os.Getenv(dbEnvName)
	if len(dbName) == 0 {
		return nil, errors.New("postgres database name not found")
	}

	username := os.Getenv(usernameEnvName)
	if len(username) == 0 {
		return nil, errors.New("postgres username not found")
	}

	password := os.Getenv(passwordEnvName)
	if len(password) == 0 {
		return nil, errors.New("postgres password not found")
	}

	return &postgresConfig{
		host,
		port,
		dbName,
		username,
		password,
	}, nil
}

func (cfg *postgresConfig) DSN() string {
	format := "host=%s port=%s dbname=%s user=%s password=%s sslmode=disable"
	return fmt.Sprintf(
		format,
		cfg.host,
		cfg.port,
		cfg.dbname,
		cfg.username,
		cfg.password,
	)
}
