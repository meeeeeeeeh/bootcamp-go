package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"log"
)

type Config struct {
	UsernameES string `envconfig:"ELASTICSEARCH_USERNAME" required:"true"`
	PasswordES string `envconfig:"ELASTICSEARCH_PASSWORD" required:"true"`
	AddressES  string `envconfig:"ELASTICSEARCH_ADDRESS" required:"true"`
	JwtSecret  string `envconfig:"JWT_SECRET" required:"true"`
	AddressSvr string `envconfig:"SERVER_ADDRESS" required:"true"`
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
