package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"log"
	"time"
)

const defaultConfigPath = "config/local.yaml"

type Config struct {
	Env        string           `yaml:"env" env:"ENV"`
	TokenTTL   time.Duration    `yaml:"token_ttl" env:"TOKEN_TTL"`
	PostgresDB PostgresDBConfig `yaml:"postgres"`
	GRPCServer GRPCServerConfig `yaml:"grpc_server"`
}

type GRPCServerConfig struct {
	Port    int           `yaml:"port" env:"GRPC_SERVER_PORT"`
	Timeout time.Duration `yaml:"timeout" env:"GRPC_SERVER_TIMEOUT"`
}

type PostgresDBConfig struct {
	Host     string `yaml:"host" env:"POSTGRES_DB_HOST"`
	Port     string `yaml:"port" env:"POSTGRES_DB_PORT"`
	DBName   string `yaml:"db_name" env:"POSTGRES_DB_NAME"`
	Username string `yaml:"username" env:"POSTGRES_DB_USERNAME"`
	Password string `yaml:"password" env:"POSTGRES_DB_PASSWORD"`
}

func MustLoad() *Config {
	_ = godotenv.Load()

	configPath := "config/local.yaml"
	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("failed to read config: %v", err)
	}

	return &cfg
}
