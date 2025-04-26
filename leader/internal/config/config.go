package config

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

const (
	EnvLocal = "local"
	EnvDev   = "dev"
	EnvProd  = "prod"
)

var cfg config

type config struct {
	EnvType  string      `yaml:"env_type"`
	Database DatabaseCfg `yaml:"POSTGRES"`
	Server   ServerCfg   `yaml:"server"`
	Logger   LoggerCfg   `yaml:"LOGGER"`
}

type ServerCfg struct {
	Host     string `yaml:"host"`
	HttpPort int    `yaml:"port"`
	Port     int    `yaml:"ROLE_LEADER_GRPC_PORT"`
}

type DatabaseCfg struct {
	Host           string `yaml:"POSTGRES_HOST"`
	Port           int    `yaml:"POSTGRES_PORT"`
	User           string `yaml:"POSTGRES_USER"`
	Password       string `yaml:"POSTGRES_PASSWORD"`
	DbName         string `yaml:"POSTGRES_DB"`
	MaxConnections int    `yaml:"POSTGRES_MAX_CONN"`
	MinConnections int    `yaml:"POSTGRES_MIN_CONN"`
}

type LoggerCfg struct {
	Env string `yaml:"ENV"`
}

// Cfg return copy of cfg (line 18)
func Config() config {
	return cfg
}

func init() {
	envType := getEnvType()
	path := getConfigFilePath(envType)
	cleanenv.ReadConfig(path, &cfg)
}

func getConfigFilePath(envType string) string {
	path := fmt.Sprintf("./config/%s.yaml", envType)
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		log.Fatalf("%s file not found", path)
	}
	return path
}

func getEnvType() string {
	envType := os.Getenv("ENV_TYPE")
	if envType == "" {
		log.Fatal("Empty ENV_TYPE variable")
	}
	if envType != EnvProd {
		log.Printf("!!! Using %s env type. Not for production !!!", envType)
		log.Printf("!!! Using %s env type. Not for production !!!", envType)
		log.Printf("!!! Using %s env type. Not for production !!!", envType)
	}
	return envType
}
