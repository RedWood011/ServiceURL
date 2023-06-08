package config

import (
	"encoding/json"
	"fmt"
	"os"
)

const path = `/config/config.json`

// Config config file
type Config struct {
	StyleCheck  []string
	StaticCheck []string
}

// ReadConfig чтение файла конфигурации
func ReadConfig(cfg *Config) error {
	dir, err := os.Getwd()
	fmt.Println(dir)
	data, err := os.ReadFile(dir + path)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(data, &cfg); err != nil {
		return err
	}
	return nil
}

// NewConfig конструктор
func NewConfig() *Config {
	return &Config{}
}
