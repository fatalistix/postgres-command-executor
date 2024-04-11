package config

import (
	"encoding/json"
	"os"
	"time"
)

type Config struct {
	HTTPServer HTTPServer `json:"http_server"`
}

type HTTPServer struct {
	Port        int           `json:"port"`
	Timeout     time.Duration `json:"timeout"`
	IdleTimeout time.Duration `json:"idle_timeout"`
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
