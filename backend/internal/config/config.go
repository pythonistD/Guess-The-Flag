package config

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"time"
)

const DefaultMaxFlagsPerGame = 30

type AppConfig struct {
	Addr             string        `yaml:"addr"`
	Secret           string        `yaml:"secret"`
	TokenLifetime    time.Duration `yaml:"token_lifetime"`
	DBConfig         DBConfig      `yaml:"database"`
	Level            string        `yaml:"logging_level"`
	MaxFlagsPerGame  int           `yaml:"max_flags_per_game"`
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
	if err := yaml.Unmarshal(data, &appConfig); err != nil {
		return AppConfig{}, fmt.Errorf("error unmarshalling config: %w", err)
	}
	if appConfig.MaxFlagsPerGame <= 0 {
		appConfig.MaxFlagsPerGame = DefaultMaxFlagsPerGame
	}
	return appConfig, nil
}
