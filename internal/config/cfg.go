package config

import (
	"flag"
	"os"
)

type Config struct {
	Port     string
	Address  string
	FilePath string
}

func NewConfig() *Config {
	cfg := &Config{}
	if cfg.Port = os.Getenv("SERVER_ADDRESS"); cfg.Port == "" {
		flag.StringVar(&cfg.Port, "a", ":8080", "Server port")
	}

	if cfg.Address = os.Getenv("BASE_URL"); cfg.Address == "" {
		flag.StringVar(&cfg.Address, "b", "http://localhost:8080/", "Server  URL")
	}

	if cfg.FilePath = os.Getenv("FILE_STORAGE_PATH"); cfg.FilePath == "" {
		flag.StringVar(&cfg.FilePath, "f", "", "File storage path")
	}

	flag.Parse()

	return cfg
}
