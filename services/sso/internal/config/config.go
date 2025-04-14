package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"log"
	"os"
	"time"
)

const defaultConfigPath = "config/local.yaml"

type Config struct {
	Env        string           `yaml:"env" env-default:"local"`
	TokenTTL   time.Duration    `yaml:"token_ttl" env-default:"1h"`
	PostgresDB PostgresDBConfig `yaml:"postgres"`
	GRPCServer GRPCServerConfig `yaml:"grpc_server"`
}

type GRPCServerConfig struct {
	Port    int           `yaml:"port" env-default:"8080"`
	Timeout time.Duration `yaml:"timeout" env-default:"1h"`
}

type PostgresDBConfig struct {
	Host     string `yaml:"host" env-default:"localhost"`
	Port     string `yaml:"port" env-default:"5432"`
	DBName   string `yaml:"db_name"`
	Username string `yaml:"username"`
	Password string `yaml:"password" env:"DB_PASSWORD"`
}

func MustLoad() *Config {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
	configPath, exists := os.LookupEnv("CONFIG_PATH")
	if !exists {
		configPath = defaultConfigPath
	}
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatal("CONFIG_PATH does not exist")
	}

	var cfg Config
	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		log.Fatal(err)
	}
	return &cfg
}
