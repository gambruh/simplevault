// Package config provides functionality to configure client and server
package config

import (
	"flag"
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
	Address     *string
	Certificate *string
	PrivateKey  *string
	Key         *string
	Database    *string
}

// UserId type is used to set server cookies
type UserID string

var (
	//global variable for server config
	Cfg Config
	//global variable for server flags
	Flags FlagConfig
	//global variable for client config
	ClientCfg ClientConfig
	//global variable for client flags
	ClientFlags ClientFlagConfig
)

// InitFlags initiates server flags, giving its default values in case if no flag is provided
func InitFlags() {
	Flags.Address = flag.String("a", "localhost:8080", "server address in format host:port")
	Flags.Certificate = flag.String("cert", "cert.pem", "certificate to run TLS")
	Flags.PrivateKey = flag.String("privatekey", "privatekey.pem", "server's private key")
	Flags.Database = flag.String("d", "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable", "postgres database uri")
	Flags.Key = flag.String("k", "abcd", "key to hash")
	flag.Parse()
}

// SetConfig looks for env values and parses flags in case if there are none
// env values preferred over flag values
func SetConfig() {
	env.Parse(&Cfg)
	if _, check := os.LookupEnv("GK_ADDRESS"); !check {
		Cfg.Address = *Flags.Address
	}
	if _, check := os.LookupEnv("GK_DATABASE"); !check {
		Cfg.Database = *Flags.Database
	}
	if _, check := os.LookupEnv("GK_CERT"); !check {
		Cfg.Certificate = *Flags.Certificate
	}
	if _, check := os.LookupEnv("GK_HASHKEY"); !check {
		Cfg.Key = *Flags.Key
	}
}
