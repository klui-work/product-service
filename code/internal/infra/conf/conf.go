package conf

import "os"

type Config struct {
	PostgresConf PostgresConf
}

type PostgresConf struct {
	PostgresHost     string
	PostgresUser     string
	PostgresPassword string
	PostgresDB       string
	PostgresPort     string
	PostgresSSLMode  string
}

type AppEnv string

var Env *Config

func LoadConfig() {
	Env = &Config{
		PostgresConf: PostgresConf{
			PostgresHost:     os.Getenv("POSTGRES_HOST"),
			PostgresUser:     os.Getenv("POSTGRES_USER"),
			PostgresPassword: os.Getenv("POSTGRES_PASSWORD"),
			PostgresDB:       os.Getenv("POSTGRES_DB"),
			PostgresPort:     os.Getenv("POSTGRES_PORT"),
			PostgresSSLMode:  os.Getenv("POSTGRES_SSLMODE"),
		},
	}
}
