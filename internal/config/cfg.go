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
	databaseDSN      = ""
	countRepetitionB = 5
	amountWorkers    = 5
	sizeBufWorker    = 100
	isHTTPS          = false
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
		fileCfg.CountRepetitionBD = countRepetitionB
		fileCfg.AmountWorkers = amountWorkers
		fileCfg.SizeBufWorker = sizeBufWorker
		fileCfg.IsHTTPS = isHTTPS
		fileCfg.KeyHash = keyHash
	}

	if cfg.ServerAddress == serverAddress {
		cfg.ServerAddress = fileCfg.ServerAddress
	}
	if cfg.BaseURL == baseURL {
		cfg.BaseURL = fileCfg.BaseURL
	}
	if cfg.FilePath == fileStoragePath {
		cfg.FilePath = fileCfg.FilePath
	}
	if cfg.DatabaseDSN == databaseDSN {
		cfg.DatabaseDSN = fileCfg.DatabaseDSN
	}
	if cfg.CountRepetitionBD == countRepetitionB {
		cfg.CountRepetitionBD = fileCfg.CountRepetitionBD
	}
	if cfg.AmountWorkers == amountWorkers {
		cfg.AmountWorkers = fileCfg.AmountWorkers
	}
	if cfg.SizeBufWorker == sizeBufWorker {
		cfg.SizeBufWorker = fileCfg.SizeBufWorker
	}
	if !cfg.IsHTTPS {
		cfg.IsHTTPS = fileCfg.IsHTTPS
	}
	if cfg.KeyHash == keyHash {
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
		flag.StringVar(&c.ServerAddress, "a", serverAddress, "Server port")
	}
	if isDefault(c.BaseURL) {
		flag.StringVar(&c.BaseURL, "b", baseURL, "Server  URL")
	}
	if isDefault(c.FilePath) {
		flag.StringVar(&c.FilePath, "f", fileStoragePath, "File storage path")
	}
	if isDefault(c.KeyHash) {
		flag.StringVar(&c.KeyHash, "key", keyHash, "KeyHash secret")
	}
	if isDefault(c.DatabaseDSN) {
		flag.StringVar(&c.DatabaseDSN, "d", databaseDSN, "Database DSN")
	}
	if isDefault(c.CountRepetitionBD) {
		flag.IntVar(&c.CountRepetitionBD, "r", countRepetitionB, "CountRepetitionBD")
	}

	if isDefault(c.AmountWorkers) {
		flag.IntVar(&c.AmountWorkers, "workers", amountWorkers, "Number of workers")
	}
	if isDefault(c.SizeBufWorker) {
		flag.IntVar(&c.SizeBufWorker, "buff", sizeBufWorker, "Workers channel buffer")
	}

	if isDefault(c.IsHTTPS) {
		flag.BoolVar(&c.IsHTTPS, "https", isHTTPS, "Enable HTTPS")
	}

	flag.StringVar(&c.ConfigPath, "config", "", "configuration file")

	flag.Parse()
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
		CountRepetitionBD: countRepetitionB,
		AmountWorkers:     amountWorkers,
		SizeBufWorker:     sizeBufWorker,
	}

}
