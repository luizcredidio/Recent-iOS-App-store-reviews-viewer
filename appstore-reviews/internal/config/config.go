package config

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type Config struct {
	Addr                string   `json:"addr"`
	PollIntervalSeconds int      `json:"pollIntervalSeconds"`
	ReviewWindowHours   int      `json:"reviewWindowHours"`
	AppIDs              []string `json:"appIDs"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	var c Config
	if err := json.Unmarshal(data, &c); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	return &c, nil
}

func (c *Config) PollInterval() time.Duration {
	return time.Duration(c.PollIntervalSeconds) * time.Second
}

func (c *Config) ReviewWindow() time.Duration {
	return time.Duration(c.ReviewWindowHours) * time.Hour
}