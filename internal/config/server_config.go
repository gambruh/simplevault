// Package config provides functionality to configure client and server
package config

import (
	"fmt"
	"os"

	"github.com/caarlos0/env/v6"
)

// Config stores server configuration
type Config struct {
	Address     string `env:"GK_ADDRESS" envDefault:"localhost:8080"`
	Certificate string `env:"GK_CERT" envDefault:"cert.pem"`
	PrivateKey  string `env:"GK_PRIVATE_KEY" envDefault:"privatekey.pem"`
	Key         string `env:"GK_HASHKEY" envDefault:""`
	Database    string `env:"GK_DATABASE" envDefault:"postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"`
}

// FlagConfig stores flag values
type FlagConfig struct {
	Address     string
	Certificate string
	PrivateKey  string
	Key         string
	Database    string
}

// UserId type is used to set server cookies
type UserID string

var (
	//global variable for server config
	Cfg Config
	//global variable for client config
	ClientCfg ClientConfig
)

// DefaultServerFlags returns default CLI values for server startup.
func DefaultServerFlags() FlagConfig {
	return FlagConfig{
		Address:     "localhost:8080",
		Certificate: "cert.pem",
		PrivateKey:  "privatekey.pem",
		Database:    "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable",
		Key:         "abcd",
	}
}

// SetConfig looks for env values and parses flags in case if there are none
// env values preferred over flag values
func SetConfig(flags FlagConfig) error {
	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		return fmt.Errorf("parse server env: %w", err)
	}
	if _, check := os.LookupEnv("GK_ADDRESS"); !check {
		cfg.Address = flags.Address
	}
	if _, check := os.LookupEnv("GK_DATABASE"); !check {
		cfg.Database = flags.Database
	}
	if _, check := os.LookupEnv("GK_CERT"); !check {
		cfg.Certificate = flags.Certificate
	}
	if _, check := os.LookupEnv("GK_PRIVATE_KEY"); !check {
		cfg.PrivateKey = flags.PrivateKey
	}
	if _, check := os.LookupEnv("GK_HASHKEY"); !check {
		cfg.Key = flags.Key
	}
	Cfg = cfg
	return nil
}
