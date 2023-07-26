package config

import (
	"flag"
	"fmt"
	"log"
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
	Address         *string
	ClientCert      *string
	ServerCert      *string
	PrivateKey      *string
	LocalStorage    *string
	UserDataFolder  *string
	UserDataFile    *string
	BinInputFolder  *string
	BinOutputFolder *string
	CheckTime       *time.Duration
}

// InitClientFlags simply initiates the client flags
func InitClientFlags() {
	ClientFlags.Address = flag.String("a", "localhost:8080", "server address in format host:port")
	ClientFlags.ClientCert = flag.String("s", "publickey.pem", "path to client's certificate file")
	ClientFlags.PrivateKey = flag.String("p", "privatekey.pem", "path to file with public key for agent")
	ClientFlags.LocalStorage = flag.String("localstorage", "./localstorage", "address of the folder to store files")
	ClientFlags.CheckTime = flag.Duration("t", 60*time.Second, "interval in time.Duration format (10s, 5m) to check data from DB")
	ClientFlags.BinInputFolder = flag.String("bininputfolder", "./filetosend", "folder to put binaries in to be sent")
	ClientFlags.BinOutputFolder = flag.String("binoutputfolder", "./filesrcv", "folder to store received binaries")
}

// SetClientConfig sets the config, parsing flags and looking for env values
// Env values are preferred over flags
func SetClientConfig() (cfg ClientConfig) {
	InitClientFlags()
	env.Parse(&ClientCfg)
	if _, check := os.LookupEnv("GK_ADDRESS"); !check {
		cfg.Address = *ClientFlags.Address
	}
	if _, check := os.LookupEnv("GK_PRIVATEKEY"); !check {
		ex, err := os.Getwd()
		if err != nil {
			log.Println("error when trying to get filepath in SetClientConfig:", err)
		}

		cfg.PrivateKey = filepath.Join(filepath.Dir(ex), *ClientFlags.PrivateKey)
		fmt.Println("CLIENT PRIVATE KEY", cfg.PrivateKey)
	}
	if _, check := os.LookupEnv("GK_LOCALSTORAGE"); !check {
		cfg.LocalStorage = *ClientFlags.LocalStorage
	}
	if _, check := os.LookupEnv("GK_CHECKINTERVAL"); !check {
		cfg.CheckTime = *ClientFlags.CheckTime
	}
	if _, check := os.LookupEnv("GK_BINARIES_INPUT"); !check {
		ClientCfg.BinInputFolder = *ClientFlags.BinInputFolder
	}
	if _, check := os.LookupEnv("GK_BINARIES_OUTPUT"); !check {
		cfg.BinOutputFolder = *ClientFlags.BinOutputFolder
	}
	if _, check := os.LookupEnv("GK_CERT"); !check {
		ex, err := os.Getwd()
		if err != nil {
			log.Println("error when trying to get filepath in SetClientConfig:", err)
		}

		cfg.ClientCert = filepath.Join(filepath.Dir(ex), *ClientFlags.ClientCert)
		fmt.Println("CLIENT CERTIFICATE", cfg.ClientCert)
	}
	return cfg
}
