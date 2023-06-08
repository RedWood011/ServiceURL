// Package config Пакет конфигурации сервиса
package config

import (
	"encoding/json"
	"flag"
	"io"
	"log"
	"os"

	"github.com/caarlos0/env"
)

const (
	serverAddress   = ":8080"
	baseURL         = "http://localhost:8080"
	fileStoragePath = ""
	keyHash         = "7cdb395a-e63e-445f-b2c4-90a400438ee4"
	//databaseDSN       = "postgres://qwerty:qwerty@localhost:5438/postgres?sslmode=disable"
	databaseDSN       = ""
	CountRepetitionBD = 5
	AmountWorkers     = 5
	SizeBufWorker     = 100
	isHTTPS           = false
)

// Config Конфигурация приложения
type Config struct {
	ServerAddress     string `env:"SERVER_ADDRESS" envDefault:""`
	BaseURL           string `env:"BASE_URL" envDefault:""`
	FilePath          string `env:"FILE_STORAGE_PATH" envDefault:""`
	KeyHash           string `env:"KEY_HASH" envDefault:""`
	DatabaseDSN       string `env:"DATABASE_DSN" envDefault:""`
	CountRepetitionBD int    `env:"REPETITION_CONNECT" envDefault:""`
	IsHTTPS           bool   `env:"ENABLE_HTTPS" envDefault:""`
	AmountWorkers     int    `env:"AMOUNT_WORKERS" envDefault:""`
	SizeBufWorker     int    `env:"BUF_WORKERS" envDefault:""`
	ConfigPath        string
}

// JSONConfig Конфигурация приложения json-файлом
type JSONConfig struct {
	BaseURL         string `json:"base_url"`
	ServerAddress   string `json:"server_address"`
	FileStoragePath string `json:"file_storage_path"`
	DatabaseDSN     string `json:"database_dsn"`
	IsHTTPS         bool   `json:"enable_https"`
}

// NewConfig Создание конфигурации
func NewConfig() *Config {

	cfg := &Config{}
	cfg.parseEnv()
	cfg.parseFlags()

	fileCfg := Config{}
	switch {
	case cfg.ConfigPath != "":
		fileCfg = readConfigFromFIle(cfg.ConfigPath)
	default:
		fileCfg.ServerAddress = serverAddress
		fileCfg.BaseURL = baseURL
		fileCfg.FilePath = fileStoragePath
		fileCfg.DatabaseDSN = databaseDSN
		fileCfg.CountRepetitionBD = CountRepetitionBD
		fileCfg.AmountWorkers = AmountWorkers
		fileCfg.SizeBufWorker = SizeBufWorker
		fileCfg.IsHTTPS = isHTTPS
		fileCfg.KeyHash = keyHash
	}

	if isDefault(cfg.ServerAddress) {
		cfg.ServerAddress = fileCfg.ServerAddress
	}
	if isDefault(cfg.BaseURL) {
		cfg.BaseURL = fileCfg.BaseURL
	}
	if isDefault(cfg.FilePath) {
		cfg.FilePath = fileCfg.FilePath
	}
	if isDefault(cfg.DatabaseDSN) {
		cfg.DatabaseDSN = fileCfg.DatabaseDSN
	}
	if isDefault(cfg.CountRepetitionBD) {
		cfg.CountRepetitionBD = fileCfg.CountRepetitionBD
	}
	if isDefault(cfg.AmountWorkers) {
		cfg.AmountWorkers = fileCfg.AmountWorkers
	}
	if isDefault(cfg.SizeBufWorker) {
		cfg.SizeBufWorker = fileCfg.SizeBufWorker
	}
	if isDefault(cfg.IsHTTPS) {
		cfg.IsHTTPS = fileCfg.IsHTTPS
	}
	if isDefault(cfg.KeyHash) {
		cfg.KeyHash = fileCfg.KeyHash
	}

	return cfg
}

func (c *Config) parseEnv() {
	err := env.Parse(c)
	if err != nil {
		log.Fatal(err)
	}
}

func (c *Config) parseFlags() {
	if isDefault(c.ServerAddress) {
		flag.StringVar(&c.ServerAddress, "a", "", "Server port")
	}
	if isDefault(c.BaseURL) {
		flag.StringVar(&c.BaseURL, "b", "", "Server  URL")
	}
	if isDefault(c.FilePath) {
		flag.StringVar(&c.FilePath, "f", "", "File storage path")
	}
	if isDefault(c.KeyHash) {
		flag.StringVar(&c.KeyHash, "key", "", "KeyHash secret")
	}
	if isDefault(c.DatabaseDSN) {
		flag.StringVar(&c.DatabaseDSN, "d", "", "Database DSN")
	}
	if isDefault(c.CountRepetitionBD) {
		flag.IntVar(&c.CountRepetitionBD, "r", 0, "CountRepetitionBD")
	}

	if isDefault(c.AmountWorkers) {
		flag.IntVar(&c.AmountWorkers, "workers", 0, "Number of workers")
	}
	if isDefault(c.SizeBufWorker) {
		flag.IntVar(&c.SizeBufWorker, "buff", 0, "Workers channel buffer")
	}

	if isDefault(c.IsHTTPS) {
		flag.BoolVar(&c.IsHTTPS, "https", false, "Enable HTTPS")
	}
	if isDefault(c.IsHTTPS) {
		flag.StringVar(&c.ConfigPath, "config", "", "configuration file")
	}
}

func isDefault[T comparable](v T) bool {
	var zero T
	return v == zero
}

func readConfigFromFIle(fileName string) Config {
	f, err := os.Open(fileName)
	if err != nil {
		log.Fatalln(err)
	}
	data, err := io.ReadAll(f)
	if err != nil {
		log.Fatal()
	}
	cfg := JSONConfig{}
	err = json.Unmarshal(data, &cfg)
	if err != nil {
		log.Fatal(err)
	}

	return Config{
		ServerAddress:     cfg.ServerAddress,
		BaseURL:           cfg.BaseURL,
		FilePath:          cfg.FileStoragePath,
		DatabaseDSN:       cfg.DatabaseDSN,
		IsHTTPS:           cfg.IsHTTPS,
		CountRepetitionBD: CountRepetitionBD,
		AmountWorkers:     AmountWorkers,
		SizeBufWorker:     SizeBufWorker,
	}

}
