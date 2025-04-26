package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	NotificationsGrpcPort int               `yaml:"NOTIFICATIONS_GRPC_PORT" env:"NOTIFICATIONS_GRPC_PORT" env-default:"50051"`
	Host                  string            `yaml:"HOST" env:"HOST" env-default:"localhost"`
	Port                  string            `yaml:"PORT" env:"PORT" env-default:"8084"`
	SendMail              ConfigForSendMail `yaml:"CREDENTIALS_SENDER" env:"SEND_MAIL"`
}

type ConfigForSendMail struct {
	SenderEmail    string `yaml:"SENDER_EMAIL" env:"SENDER_EMAIL"`
	SenderPassword string `yaml:"SENDER_PASSWORD" env:"SENDER_PASSWORD"`
}

func New() (*Config, error) {
	var cfg Config

	if err := cleanenv.ReadConfig("./config/config.yaml", &cfg); err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	return &cfg, nil
}
