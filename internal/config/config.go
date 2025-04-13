package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	LogConfig    LogConfig
	ServerConfig ServerConfig
	DBConfig     DBConfig
}

type LogConfig struct {
	LogLevel string `env:"LOG_LEVEL" envDefault:"INFO"`
}

type ServerConfig struct {
	ServerHTTPPort  int    `env:"SERVER_HTTP_PORT,required"`
	ServerSecretKey string `env:"SERVER_SECRET_KEY,required"`
}

type DBConfig struct {
	PostgresDSN string `env:"POSTGRES_DSN,required"`

	GooseDriver       string `env:"GOOSE_DRIVER" envDefault:"postgres"`
	GooseMigrationDir string `env:"GOOSE_MIGRATION_DIR" envDefault:"migrations"`
}

func NewConfig() (*Config, error) {
	var cfg Config

	if err := env.Parse(&cfg); err != nil {
		return nil, fmt.Errorf("configuration parsing: %w", err)
	}

	return &cfg, nil
}
