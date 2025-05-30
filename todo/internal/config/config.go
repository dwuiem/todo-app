package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"log"
	"time"
)

type Config struct {
	Env         string        `yaml:"env" env:"ENV"`
	StoragePath string        `yaml:"storage_path" env:"STORAGE_PATH"`
	AppSecret   string        `yaml:"app_secret" env:"APP_SECRET"`
	AppId       int           `yaml:"app_id" env:"APP_ID"`
	HTTPServer  HTTPServer    `yaml:"http_server"`
	PostgresDB  PostgresDB    `yaml:"postgres_db"`
	Clients     ClientsConfig `yaml:"clients"`
}

type HTTPServer struct {
	Addr        string        `yaml:"address" env:"HTTP_SERVER_ADDRESS"`
	Timeout     time.Duration `yaml:"timeout" env:"HTTP_SERVER_TIMEOUT"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env:"HTTP_SERVER_IDLE_TIMEOUT"`
}

type PostgresDB struct {
	Host     string `yaml:"host" env:"POSTGRES_DB_HOST"`
	Port     string `yaml:"port" env:"POSTGRES_DB_PORT"`
	DBName   string `yaml:"db_name" env:"POSTGRES_DB_NAME"`
	Username string `yaml:"username" env:"POSTGRES_DB_USERNAME"`
	Password string `yaml:"password" env:"POSTGRES_DB_PASSWORD"`
}

type ClientsConfig struct {
	SSO ClientConfig `yaml:"sso"`
}

type ClientConfig struct {
	Address      string        `yaml:"address" env:"CLIENT_SSO_ADDRESS"`
	Timeout      time.Duration `yaml:"timeout" env:"CLIENT_SSO_TIMEOUT"`
	RetriesCount int           `yaml:"retries_count" env:"CLIENT_SSO_RETRIES_COUNT"`
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
