package config

import (
	"encoding/json"
	"os"
	"time"
)

type Config struct {
	HTTPServer HTTPServer `json:"http_server"`
	Postgres   Postgres   `json:"postgres"`
}

type HTTPServer struct {
	Address         string        `json:"address"`
	Timeout         time.Duration `json:"timeout"`
	IdleTimeout     time.Duration `json:"idle_timeout"`
	ShutdownTimeout time.Duration `json:"shutdown_timeout"`
}

type Postgres struct {
	Host     string `json:"host"`
	Port     uint16 `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	DBName   string `json:"dbname"`
	SSLMode  string `json:"sslmode"`
}

func MustLoadConfig() Config {
	pathToConfig := os.Getenv("CONFIG_PATH")

	data, err := os.ReadFile(pathToConfig)
	if err != nil {
		panic("config file not found: " + pathToConfig)
	}

	var config Config
	if err = json.Unmarshal(data, &config); err != nil {
		panic("cannot read config: " + err.Error())
	}

	return config
}
