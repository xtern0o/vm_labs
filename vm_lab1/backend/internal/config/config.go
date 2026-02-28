package config

import (
	"errors"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Solver SolverConfig `yaml:"solver"`
}

type SolverConfig struct {
	Epsilon float64 `yaml:"epsilon"`
	MaxIter int     `yaml:"max_iter"`
}

func MustLoad(path string) (Config, error) {
	if path == "" {
		return Config{}, errors.New("config path is empty")
	}

	b, err := os.ReadFile(path)
	if err != nil {
		return Config{}, fmt.Errorf("read config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(b, &cfg); err != nil {
		return Config{}, fmt.Errorf("parse error: %w", err)
	}

	return cfg, nil

}
