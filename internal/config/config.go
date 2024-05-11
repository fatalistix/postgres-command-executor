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
	Address         string        `json:"address"`
	Timeout         time.Duration `json:"timeout"`
	IdleTimeout     time.Duration `json:"idle_timeout"`
	ShutdownTimeout time.Duration `json:"shutdown_timeout"`
}

func MustLoadConfig(configPath string) Config {
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("config file does not exists: " + configPath)
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		panic("cannot read config file: " + err.Error())
	}

	var config Config
	if err = json.Unmarshal(data, &config); err != nil {
		panic("cannot unmarshal config to json: " + err.Error())
	}

	return config
}
