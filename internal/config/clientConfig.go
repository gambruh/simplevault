package config

import (
	"flag"
	"os"
	"time"

	"github.com/caarlos0/env/v6"
)

type ClientConfig struct {
	Address         string        `env:"GK_ADDRESS" envDefault:"localhost:8080"`
	PublicKey       string        `env:"GK_PUBLICKEY" envDefault:"publickey.pem"`
	LocalStorage    string        `env:"GK_LOCALSTORAGE" envDefault:"./localstorage"`
	BinInputFolder  string        `env:"GK_BINARIES_INPUT" envDefault:"./filetosend"`
	UserData        string        `env:"GK_USERDATA" envDefault:"./userdata/user.json"`
	BinOutputFolder string        `env:"GK_BINARIES_OUTPUT" envDefault:"./filesrcv"`
	CheckTime       time.Duration `env:"GK_CHECKINTERVAL" envDefault:"60s"`
}

type ClientFlagConfig struct {
	Address         *string
	PublicKey       *string
	LocalStorage    *string
	UserData        *string
	BinInputFolder  *string
	BinOutputFolder *string
	CheckTime       *time.Duration
}

func InitClientFlags() {
	ClientFlags.Address = flag.String("a", "localhost:8080", "server address in format host:port")
	ClientFlags.PublicKey = flag.String("p", "publickey.pem", "path to file with public key for agent")
	ClientFlags.LocalStorage = flag.String("localstorage", "./localstorage", "address of the folder to store files")
	ClientFlags.CheckTime = flag.Duration("t", 60*time.Second, "interval in time.Duration format (10s, 5m) to check data from DB")
	ClientFlags.BinInputFolder = flag.String("bininputfolder", "./filetosend", "folder to put binaries in to be sent")
	ClientFlags.BinOutputFolder = flag.String("binoutputfolder", "./filesrcv", "folder to store received binaries")
}

func SetClientConfig() {
	env.Parse(&ClientCfg)
	if _, check := os.LookupEnv("GK_ADDRESS"); !check {
		ClientCfg.Address = *ClientFlags.Address
	}
	if _, check := os.LookupEnv("GK_PUBLICKEY"); !check {
		ClientCfg.PublicKey = *ClientFlags.PublicKey
	}
	if _, check := os.LookupEnv("GK_LOCALSTORAGE"); !check {
		ClientCfg.LocalStorage = *ClientFlags.LocalStorage
	}
	if _, check := os.LookupEnv("GK_CHECKINTERVAL"); !check {
		ClientCfg.CheckTime = *ClientFlags.CheckTime
	}
	if _, check := os.LookupEnv("GK_BINARIES_INPUT"); !check {
		ClientCfg.BinInputFolder = *ClientFlags.BinInputFolder
	}
	if _, check := os.LookupEnv("GK_BINARIES_OUTPUT"); !check {
		ClientCfg.BinOutputFolder = *ClientFlags.BinOutputFolder
	}
}
