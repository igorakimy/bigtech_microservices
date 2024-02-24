package config

import (
	"github.com/joho/godotenv"
)

type GRPCConfig interface {
	Address() string
}

type PostgresConfig interface {
	DSN() string
}

func Load(path string) error {
	if err := godotenv.Load(path); err != nil {
		return err
	}

	return nil
}
