// Package config Пакет конфигурации сервиса
package config

import (
	"flag"
	"os"
)

// Config Конфигурация приложения
type Config struct {
	Port              string
	Address           string
	FilePath          string
	KeyHash           string
	DatabaseDSN       string
	CountRepetitionBD string
	NumWorkers        int
	SizeBufWorker     int
}

// NewConfig Создание конфигурации
func NewConfig() *Config {
	cfg := &Config{}
	if cfg.Port = os.Getenv("SERVER_ADDRESS"); cfg.Port == "" {
		flag.StringVar(&cfg.Port, "a", ":8080", "Server port")
	}

	if cfg.Address = os.Getenv("BASE_URL"); cfg.Address == "" {
		flag.StringVar(&cfg.Address, "b", "http://localhost:8080", "Server  URL")
	}

	if cfg.FilePath = os.Getenv("FILE_STORAGE_PATH"); cfg.FilePath == "" {
		//flag.StringVar(&cfg.FilePath, "f", "/Users/evyaroshen/GolandProjects/yandex/ServiceURL/words.json", "File storage path")
		flag.StringVar(&cfg.FilePath, "f", "", "File storage path")
	}

	if cfg.KeyHash = os.Getenv("KEY_HASH"); cfg.KeyHash == "" {
		flag.StringVar(&cfg.KeyHash, "key", "7cdb395a-e63e-445f-b2c4-90a400438ee4", "KeyHash secret")
	}

	if cfg.DatabaseDSN = os.Getenv("DATABASE_DSN"); cfg.DatabaseDSN == "" {
		flag.StringVar(&cfg.DatabaseDSN, "d", "", "Database connection")
		//flag.StringVar(&cfg.DatabaseDSN, "d", "postgres://qwerty:qwerty@localhost:5438/postgres?sslmode=disable", "")
	}

	if cfg.CountRepetitionBD = os.Getenv("REPETITION_CONNECT"); cfg.CountRepetitionBD == "" {
		flag.StringVar(&cfg.CountRepetitionBD, "repetition", "5", "repetition connect database")
	}
	cfg.NumWorkers = 5
	cfg.SizeBufWorker = 100

	flag.Parse()

	return cfg
}
