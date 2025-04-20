package config

import "github.com/ilyakaznacheev/cleanenv"

type Config struct {
	HttpServer HttpServer
	Database   Database
	JWT        JWT
}

type HttpServer struct {
	Port string `env:"PORT" env-default:"8080"`
}

type Database struct {
	Host     string `env:"DB_HOST" env-default:"localhost"`
	Port     string `env:"DB_PORT" env-default:"5432"`
	User     string `env:"DB_USER" env-default:"postgres"`
	Password string `env:"DB_PASSWORD" env-required:"true"`
	Name     string `env:"DB_NAME" env-default:"postgres"`
}

type JWT struct {
	Secret string `env:"JWT_SECRET" env-required:"true"`
}

func MustInit(path string) *Config {
	var cfg Config

	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic(err)
	}

	return &cfg
}
