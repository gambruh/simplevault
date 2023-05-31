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
	//	Storage     bool
}

type ClientConfig struct {
	Address   string `env:"GK_ADDRESS" envDefault:"localhost:8080"`
	PublicKey string `env:"GK_PUBLICKEY" envDefault:"publickey.pem"`
	Key       string `env:"GK_HASHKEY" envDefault:""`
}

type FlagConfig struct {
	Address     *string
	Certificate *string
	Key         *string
	Database    *string
	Storage     *bool
}

type ClientFlagConfig struct {
	User      *string
	Address   *string
	PublicKey *string
	Key       *string
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
	//Flags.Storage = flag.Bool("s", false, "inmemory storage for lazy debugging")
	flag.Parse()
}

func InitClientFlags() {
	ClientFlags.Address = flag.String("a", "localhost:8080", "server address in format host:port")
	ClientFlags.PublicKey = flag.String("p", "publickey.pem", "path to file with public key for agent")
	ClientFlags.Key = flag.String("k", "abcd", "key to hash files")
}

func SetClientConfig() {
	env.Parse(&ClientCfg)
	if _, check := os.LookupEnv("GK_ADDRESS"); !check {
		ClientCfg.Address = *ClientFlags.Address
	}
	if _, check := os.LookupEnv("GK_PUBLICKEY"); !check {
		ClientCfg.PublicKey = *ClientFlags.PublicKey
	}
	if _, check := os.LookupEnv("GK_HASHKEY"); !check {
		ClientCfg.Key = *ClientFlags.Key
	}
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

	//Cfg.Storage = *Flags.Storage
}
