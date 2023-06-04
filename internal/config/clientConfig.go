package config

import (
	"flag"
	"os"

	"github.com/caarlos0/env/v6"
)

type ClientConfig struct {
	Address      string `env:"GK_ADDRESS" envDefault:"localhost:8080"`
	PublicKey    string `env:"GK_PUBLICKEY" envDefault:"publickey.pem"`
	Key          string `env:"GK_HASHKEY" envDefault:""`
	LocalStorage string `env:"LOCALSTORAGE" envDefault:"./localstorage"`
}

type ClientFlagConfig struct {
	Address      *string
	PublicKey    *string
	Key          *string
	LocalStorage *string
}

func InitClientFlags() {
	ClientFlags.Address = flag.String("a", "localhost:8080", "server address in format host:port")
	ClientFlags.PublicKey = flag.String("p", "publickey.pem", "path to file with public key for agent")
	ClientFlags.Key = flag.String("k", "abcd", "key to hash files")
	ClientFlags.LocalStorage = flag.String("localstorage", "./localstorage", "address of the folder to store files")
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
	if _, check := os.LookupEnv("LOCALSTORAGE"); !check {
		ClientCfg.LocalStorage = *ClientFlags.LocalStorage
	}
}
