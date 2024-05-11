package env

import (
	"github.com/joho/godotenv"
	"os"
	"strconv"
)

type Env struct {
	ConfigPath       string
	PostgresHost     string
	PostgresPort     uint16
	PostgresUser     string
	PostgresPassword string
	PostgresSSLMode  string
	PostgresDB       string
}

func MustLoadEnv() *Env {
	err := godotenv.Load()
	if err != nil {
		panic("cannot load environment variables: " + err.Error())
	}

	configPath := mustGetEnv("CONFIG_PATH")
	postgresHost := mustGetEnv("POSTGRES_HOST")
	postgresPort := mustGetEnv("POSTGRES_PORT")
	postgresUser := mustGetEnv("POSTGRES_USER")
	postgresPassword := mustGetEnv("POSTGRES_PASSWORD")
	postgresSSLMode := mustGetEnv("POSTGRES_SSL_MODE")
	postgresDb := mustGetEnv("POSTGRES_DB")

	postgresPortUint, err := strconv.ParseUint(postgresPort, 10, 16)
	if err != nil {
		panic("cannot parse POSTGRES_PORT: " + err.Error())
	}

	return &Env{
		ConfigPath:       configPath,
		PostgresHost:     postgresHost,
		PostgresPort:     uint16(postgresPortUint),
		PostgresUser:     postgresUser,
		PostgresPassword: postgresPassword,
		PostgresSSLMode:  postgresSSLMode,
		PostgresDB:       postgresDb,
	}
}

func mustGetEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		panic(key + " is not set")
	}
	return value
}
