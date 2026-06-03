package config

import (
	"fmt"
	"log"
	"os"

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
	Host     string `env:"DB_HOST" env-default:"127.0.0.1"`
	Port     string `env:"DB_PORT" env-required:"true"`
	User     string `env:"DB_USER" env-required:"true"`
	Password string `env:"DB_PASSWORD" env-required:"true"`
	Name     string `env:"DB_NAME" env-required:"true"`
}
type Redis struct {
	Addr     string `env:"REDIS_ADDR" env-default:"127.0.0.1:6379"`
	Password string `env:"REDIS_PASSWORD" env-required:"true"`
}

func (c *Config) DatabaseURL() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		c.DB.User, c.DB.Password, c.DB.Host, c.DB.Port, c.DB.Name)
}
func MustLoad() *Config {
	var cfg Config

	if _, err := os.Stat(".env"); err == nil {
		if err := cleanenv.ReadConfig(".env", &cfg); err != nil {
			log.Fatalf("Критическая ошибка при чтении .env файла: %v", err)
		}
	} else if os.IsNotExist(err) {
		if err := cleanenv.ReadEnv(&cfg); err != nil {
			log.Fatalf("Критическая ошибка при чтении переменных окружения: %v", err)
		}
	} else {
		log.Fatalf("Ошибка при проверке наличия .env файла: %v", err)
	}

	return &cfg
}
