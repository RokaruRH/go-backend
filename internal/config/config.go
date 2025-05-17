package config

import (
	"log"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env         string `yaml:"environment" env-required:"true"`
	StoragePath string `yaml:"storage_path" env-required:"true"`
	HTTPServer  `yaml:"http_server" env-required:"true"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-required:"true"`
	Timeout     time.Duration `yaml:"timeout" env-required:"true"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
	User        string        `yaml: "user"     env-required:"true" `
	Password    string        `yaml: "password" env-required:"true" env: "HTTP_SERVER_PASSWORD"`
}

const CONFIG_PATH = "config/local.yaml"

func MustLoad() *Config {
	var config Config
	if err := cleanenv.ReadConfig(CONFIG_PATH, &config); err != nil {
		log.Fatalf("config file does not exist: %s", err)
	}
	return &config
}
