package config

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/caarlos0/env/v6"
)

// ClientConfig is a structure to store client configuration
type ClientConfig struct {
	Address         string        `env:"GK_ADDRESS" envDefault:"localhost:8080"`
	ClientCert      string        `env:"GK_PUBLICKEY" envDefault:"publickey.pem"`
	ServerCert      string        `env:"GK_CERT" envDefault:"cert.pem"`
	PrivateKey      string        `env:"GK_PRIVATEKEY" envDefault:"privatekey.pem"`
	LocalStorage    string        `env:"GK_LOCALSTORAGE" envDefault:"./localstorage"`
	BinInputFolder  string        `env:"GK_BINARIES_INPUT" envDefault:"./filetosend"`
	UserDataFolder  string        `env:"GK_USERDATA_FOLDER" envDefault:"./userdata"`
	UserDataFile    string        `env:"GK_USERDATA_FILE" envDefault:"./userdata/user.json"`
	BinOutputFolder string        `env:"GK_BINARIES_OUTPUT" envDefault:"./filesrcv"`
	CheckTime       time.Duration `env:"GK_CHECKINTERVAL" envDefault:"60s"`
}

// ClientFlagConfig is a structure to store client flag values
type ClientFlagConfig struct {
	Address         string
	ClientCert      string
	PrivateKey      string
	LocalStorage    string
	BinInputFolder  string
	BinOutputFolder string
	CheckTime       time.Duration
}

// DefaultClientFlags returns default CLI values for client startup.
func DefaultClientFlags() ClientFlagConfig {
	return ClientFlagConfig{
		Address:         "localhost:8080",
		ClientCert:      "publickey.pem",
		PrivateKey:      "privatekey.pem",
		LocalStorage:    "./localstorage",
		CheckTime:       60 * time.Second,
		BinInputFolder:  "./filetosend",
		BinOutputFolder: "./filesrcv",
	}
}

// SetClientConfig sets the config, parsing flags and looking for env values
// Env values are preferred over flags
func SetClientConfig(flags ClientFlagConfig) (ClientConfig, error) {
	cfg := ClientConfig{}
	if err := env.Parse(&cfg); err != nil {
		return ClientConfig{}, fmt.Errorf("parse client env: %w", err)
	}

	if _, check := os.LookupEnv("GK_ADDRESS"); !check {
		cfg.Address = flags.Address
	}
	if _, check := os.LookupEnv("GK_PRIVATEKEY"); !check {
		privateKeyPath, err := resolveClientPath(flags.PrivateKey)
		if err != nil {
			return ClientConfig{}, fmt.Errorf("resolve private key path: %w", err)
		}
		cfg.PrivateKey = privateKeyPath
	}
	if _, check := os.LookupEnv("GK_LOCALSTORAGE"); !check {
		cfg.LocalStorage = flags.LocalStorage
	}
	if _, check := os.LookupEnv("GK_CHECKINTERVAL"); !check {
		cfg.CheckTime = flags.CheckTime
	}
	if _, check := os.LookupEnv("GK_BINARIES_INPUT"); !check {
		cfg.BinInputFolder = flags.BinInputFolder
	}
	if _, check := os.LookupEnv("GK_BINARIES_OUTPUT"); !check {
		cfg.BinOutputFolder = flags.BinOutputFolder
	}
	if _, check := os.LookupEnv("GK_PUBLICKEY"); !check {
		clientCertPath, err := resolveClientPath(flags.ClientCert)
		if err != nil {
			return ClientConfig{}, fmt.Errorf("resolve client cert path: %w", err)
		}
		cfg.ClientCert = clientCertPath
	}

	ClientCfg = cfg
	return cfg, nil
}

func resolveClientPath(path string) (string, error) {
	ex, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return filepath.Join(filepath.Dir(ex), path), nil
}
