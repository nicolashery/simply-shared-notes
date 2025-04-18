package app

import (
	"fmt"
	"os"
	"strings"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	Port        int    `env:"PORT" envDefault:"3000"`
	IsDev       bool   `env:"DEV" envDefault:"false"`
	DatabaseURL string `env:"DATABASE_URL" envDefault:"sqlite:data/app.sqlite"`
}

func NewConfig() (Config, error) {
	err := godotenv.Load()
	if err != nil && !os.IsNotExist(err) {
		return Config{}, fmt.Errorf("failed to load .env file: %w", err)
	}

	var cfg Config
	err = env.Parse(&cfg)
	if err != nil {
		return Config{}, fmt.Errorf("failed to parse config from environment variables: %w", err)
	}

	return cfg, nil
}

func (c *Config) DatabasePath() string {
	return strings.TrimPrefix(c.DatabaseURL, "sqlite:")
}
