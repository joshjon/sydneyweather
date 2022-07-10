package config

import (
	"errors"
	"flag"
	"os"
	"time"

	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v2"
)

type Config struct {
	ServerPort         int           `yaml:"serverPort" validate:"required"`
	CacheExpiry        time.Duration `yaml:"cacheExpiry" validate:"required"`
	WeatherStackAPIKey string        `yaml:"weatherStackAPIKey" envconfig:"WEATHER_STACK_KEY" validate:"required"`
	OpenWeatherAPIKey  string        `yaml:"openWeatherAPIKey" envconfig:"OPEN_WEATHER_KEY" validate:"required"`
}

func Load() (*Config, error) {
	var cfg Config

	var filepath string
	flag.StringVar(&filepath, "config", "", "yaml config file path")
	flag.Parse()
	if filepath == "" {
		return nil, errors.New("required flag 'config' was not provided")
	}

	if err := readFile(filepath, &cfg); err != nil {
		return nil, err
	}

	if err := readEnv(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func readFile(filePath string, out *Config) error {
	f, err := os.Open(filePath)
	if err != nil {
		return err
	}

	if err = yaml.NewDecoder(f).Decode(out); err != nil {
		return err
	}

	return f.Close()
}

func readEnv(out *Config) error {
	err := envconfig.Process("", out)
	if err != nil {
		return err
	}

	return nil
}
