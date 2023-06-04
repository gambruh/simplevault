package config

import (
	"flag"
	"os"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	Address     string `env:"GK_ADDRESS" envDefault:"localhost:8080"`
	Certificate string `env:"GK_CERT" envDefault:"cert.pem"`
	Key         string `env:"GK_HASHKEY" envDefault:""`
	Database    string `env:"GK_DATABASE" envDefault:"postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"`
}

type FlagConfig struct {
	Address     *string
	Certificate *string
	Key         *string
	Database    *string
	Storage     *bool
}

type UserID string

var (
	Cfg         Config
	Flags       FlagConfig
	ClientCfg   ClientConfig
	ClientFlags ClientFlagConfig
)

func InitFlags() {
	Flags.Address = flag.String("a", "localhost:8080", "server address in format host:port")
	Flags.Certificate = flag.String("c", "cert.pem", "certificate to run TLS")
	Flags.Database = flag.String("d", "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable", "postgres database uri")
	Flags.Key = flag.String("k", "abcd", "key to hash")
	flag.Parse()
}

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
