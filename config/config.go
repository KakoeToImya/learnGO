package config

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env      string `env:"APP_ENV" env-default:"development"`
	Host     string `env:"APP_HOST" env-default:"0.0.0.0"`
	Port     string `env:"APP_PORT" env-default:"8080"`
	Key      string `env:"JWT_SECRET" env-required:"true"`
	LogLevel string `env:"LOG_LEVEL" env-default:"debug"`

	DB    Database
	Redis Redis
}

type Database struct {
	Host     string `env:"DB_HOST" env-default:"127.0.0.1:"`
	Port     string `env:"DB_PORT" env-required:"true"`
	User     string `env:"DB_USER" env-required:"true"`
	Password string `env:"DB_PASSWORD" env-required:"true"`
	Name     string `env:"DB_NAME" env-required:"true"`
}
type Redis struct {
	Addr     string `env:"REDIS_ADDR" env-default:"127.0.0.1:6379"`
	Password string `env:"REDIS_PASSWORD" env-required:"true"`
}

func MustLoad() *Config {
	var cfg Config

	err := cleanenv.ReadConfig(".env", &cfg)
	if err != nil {
		log.Fatalf("Критическая ошибка: не удалось загрузить конфиг: %v", err)
	}

	return &cfg
}
