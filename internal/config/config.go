package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"log"
	"os"
	"time"
)

type Config struct {
	Env         string     `yaml:"env"`
	StoragePath string     `yaml:"storage_path"`
	HTTPServer  HTTPServer `yaml:"http_server"`
	PostgresDB  PostgresDB `yaml:"postgres_db"`
}

type HTTPServer struct {
	Addr        string        `yaml:"address" env-default:"localhost:8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

type PostgresDB struct {
	Host     string `yaml:"host" env-default:"localhost"`
	Port     string `yaml:"port" env-default:"5432"`
	DBName   string `yaml:"db_name"`
	Username string `yaml:"username"`
	Password string `env:"DB_PASSWORD"`
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func MustLoad() *Config {
	configPath, exists := os.LookupEnv("CONFIG_PATH")
	if !exists {
		log.Fatal("CONFIG_PATH environment variable not set")
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
