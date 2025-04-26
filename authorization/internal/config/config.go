package config

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/ilyakaznacheev/cleanenv"
)

const (
	EnvLocal = "local"
	EnvDev   = "dev"
	EnvProd  = "prod"
)

var cfg config

type config struct {
	EnvType string        `yaml:"env_type"`
	JWT     jwtConfig     `yaml:"jwt"`
	Storage storageConfig `yaml:"storage"`
	Server  serverConfig  `yaml:"server"`
}

type jwtConfig struct {
	TokenTTl       time.Duration `yaml:"token_ttl"`
	PathPrivateKey string        `yaml:"private_key_path"` // Ignore this field
	PathPublicKey  string        `yaml:"public_key_path"`  // Ignore this field
	PrivateKey     *rsa.PrivateKey
	PublicKey      *rsa.PublicKey
}

type serverConfig struct {
	Port int `yaml:"port"`
}

type storageConfig struct {
	Host   string `yaml:"host"`
	Port   int    `yaml:"port"`
	User   string `yaml:"user"`
	Pass   string `yaml:"pass"`
	DbName string `yaml:"db_name"`
}

// Cfg return copy of cfg (line 21)
func Config() config {
	return cfg
}

func LoadConfig() {
	envType := getEnvType()
	path := getConfigFilePath(envType)
	cleanenv.ReadConfig(path, &cfg)
	readKeys()
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
		// log.Fatal("Empty ENV_TYPE variable")
		envType = EnvLocal
	}
	if envType != EnvProd {
		log.Printf("!!! Using %s env type. Not for production !!!", envType)
		log.Printf("!!! Using %s env type. Not for production !!!", envType)
		log.Printf("!!! Using %s env type. Not for production !!!", envType)
	}
	return envType
}

func readKeys() {
	// Read private
	data, err := os.ReadFile(cfg.JWT.PathPrivateKey)
	if err != nil {
		log.Fatal(err, cfg.JWT.PathPrivateKey)
	}
	cfg.JWT.PrivateKey, err = jwt.ParseRSAPrivateKeyFromPEM(data)
	if err != nil {
		log.Fatal(err)
	}

	// Read public
	data, err = os.ReadFile(cfg.JWT.PathPublicKey)
	if err != nil {
		log.Fatal(err, cfg.JWT.PathPublicKey)
	}
	cfg.JWT.PublicKey, err = jwt.ParseRSAPublicKeyFromPEM(data)
	if err != nil {
		log.Fatal(err)
	}
}
