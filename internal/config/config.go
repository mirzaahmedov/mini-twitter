package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Postgres struct {
		User    string `env:"PG_USER" env-default:"postgres"`
		PWD     string `env:"PG_PWD"`
		Host    string `env:"PG_HOST" env-default:"127.0.0.1"`
		Port    string `env:"PG_PORT" env-default:"5432"`
		DB      string `env:"PG_DB" env-default:"twitter"`
		SSLMode string `env:"PG_SSL_MODE" env-default:"disable"`
	}
	HTTP struct {
		Addr string `env:"HTTP_ADDR" env-default:":8080"`
	}
	JWT struct {
		Secret string `env:"JWT_SECRET"`
	}
}

func Load() (*Config, error) {
	c := new(Config)

	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("error loading .env file %v", err)
	}

	err = cleanenv.ReadEnv(c)
	if err != nil {
		return nil, fmt.Errorf("error loading env variables %v", err)
	}

	return c, nil
}
