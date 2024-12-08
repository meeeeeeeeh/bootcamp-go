package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"log"
)

type Config struct {
	Postgres Postgres
	Service  Service
}

type Postgres struct {
	Host     string `envconfig:"POSTGRES_HOST" required:"true"`
	Port     string `envconfig:"POSTGRES_PORT" required:"true"`
	User     string `envconfig:"POSTGRES_USER" required:"true"`
	Password string `envconfig:"POSTGRES_PASSWORD" required:"true"`
	Name     string `envconfig:"POSTGRES_NAME" required:"true"`
}

type Service struct {
	Host string `envconfig:"SERVICE_HOST" required:"true"`
	Port string `envconfig:"SERVICE_PORT" required:"true"`
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load(".env")
	if err != nil {
		return nil, fmt.Errorf("cannot load config data, err: %v", err)
	}

	cfg := Config{}
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, fmt.Errorf("cannot process config data, err: %v", err)
	} else {
		log.Println("config initialized")
	}

	return &cfg, nil
}
