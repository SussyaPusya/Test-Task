package config

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Postgres
	Rest
	Logger
	ExternalAPI
}

type ExternalAPI struct {
	Age         string `env:"AGE_API_URL"`
	Gender      string `env:"GENDER_API_URL"`
	Nationality string `env:"NATIONALITY_API_URL"`
}

type Postgres struct {
	Host     string `env:"PG_HOST"`
	Port     int    `env:"PG_PORT"`
	Database string `env:"PG_DATABASE_DATA"`
	User     string `env:"PG_USER"`
	Password string `env:"PG_PASS"`
	MaxConn  int32  `env:"PG_MAXCONN"`
	MinConn  int32  `env:"PG_MINCONN"`
}

type Rest struct {
	Port int `env:"REST_PORT_INP"`
}
type Logger struct {
	Level string `env:"LOG_LEVEL"`
}

func NewConfig() (*Config, error) {
	var cfg Config

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		log.Println(err)
	}

	cfg.Postgres.Host = "localhost"
	cfg.Postgres.Port = 5438
	cfg.Postgres.Database = "postgres"
	cfg.Postgres.User = "postgres"
	cfg.Postgres.Password = "1234"
	cfg.Postgres.MaxConn = 10
	cfg.Postgres.MinConn = 5
	cfg.Rest.Port = 8080

	return &cfg, nil
}
