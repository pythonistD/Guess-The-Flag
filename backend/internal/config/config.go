package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"
)

type AppConfig struct {
	Addr     string   `yaml:"addr"`
	DBConfig DBConfig `yaml:"dbConfig"`
}

type DBConfig struct {
	Host                  string        `yaml:"host"`
	Port                  string        `yaml:"port"`
	Username              string        `yaml:"username"`
	Password              string        `yaml:"password"`
	DbName                string        `yaml:"dbname"`
	MaxIdleConnections    int           `yaml:"max_idle_connections"`
	MaxOpenConnections    int           `yaml:"max_open_connections"`
	MaxConnectionLifeTime time.Duration `yaml:"max_connection_lifetime"`
}

func LoadConfigFromFile(path string) (AppConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return AppConfig{}, fmt.Errorf("error reading config file: %w", err)
	}
	if err != nil && errors.Is(err, os.ErrNotExist) {
		return AppConfig{}, fmt.Errorf("error file does not exist: %w", err)
	}
	return loadFromBytes(data)
}

func loadFromBytes(data []byte) (AppConfig, error) {
	var appConfig AppConfig
	if err := json.Unmarshal(data, &appConfig); err != nil {
		return AppConfig{}, fmt.Errorf("error unmarshalling config: %w", err)
	}
	return appConfig, nil
}
