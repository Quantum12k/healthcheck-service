package config

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v3"

	"github.com/Quantum12k/healthcheck-service/internal/healthcheck"
	"github.com/Quantum12k/healthcheck-service/internal/logger"
	"github.com/Quantum12k/healthcheck-service/internal/postgresql"
)

type (
	Config struct {
		Logger           logger.Config     `yaml:"logger"`
		PostgreSQL       postgresql.Config `yaml:"db"`
		CheckIntervalSec int               `yaml:"check_interval_sec"`
		URLs             []healthcheck.URL `yaml:"urls"`
	}
)

func New(settingsPath string, urlsPath string) (*Config, error) {
	cfg := &Config{}

	bytes, err := ioutil.ReadFile(settingsPath)
	if err != nil {
		return nil, fmt.Errorf("read file from path '%s': %v", settingsPath, err)
	}

	if err = yaml.Unmarshal(bytes, &cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config: %v", err)
	}

	bytes, err = ioutil.ReadFile(urlsPath)
	if err != nil {
		return nil, fmt.Errorf("read file from path '%s': %v", urlsPath, err)
	}

	if err = yaml.Unmarshal(bytes, &cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config: %v", err)
	}

	return cfg, nil
}
