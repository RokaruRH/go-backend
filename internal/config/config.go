package config

type Config struct {
	Env string `yaml:"environment" env:"ENV" envDefault:"development"`
}
