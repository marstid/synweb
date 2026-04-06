package config

import (
	"os"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	SyntheticAPIKey string `envconfig:"SYNTHETIC_API_KEY" required:"true"`
	LogLevel        string `envconfig:"LOG_LEVEL" default:"info"`
	APIBaseURL      string `envconfig:"API_BASE_URL" default:"https://api.synthetic.new/v2"`
}

func Load() (*Config, error) {
	var cfg Config

	err := envconfig.Process("synweb", &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

func Validate(cfg *Config) error {
	if cfg.SyntheticAPIKey == "" {
		return ErrMissingAPIKey
	}

	return nil
}

var ErrMissingAPIKey = &ConfigError{
	Message: "SYNTHETIC_API_KEY environment variable is required",
	Code:    "MISSING_API_KEY",
}

type ConfigError struct {
	Message string
	Code    string
}

func (e *ConfigError) Error() string {
	return e.Message
}

func init() {
	if os.Getenv("SYNTHETIC_API_KEY") == "" {
		printMissingAPIKeyError()
	}
}

func printMissingAPIKeyError() {
	println("ERROR: SYNTHETIC_API_KEY environment variable is required")
	os.Exit(1)
}
