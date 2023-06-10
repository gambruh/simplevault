package clientfunc

import (
	"fmt"
	"net/http"

	"github.com/gambruh/gophkeeper/internal/config"
	"github.com/gambruh/gophkeeper/internal/database"
	"github.com/gambruh/gophkeeper/internal/localstorage"
)

type Client struct {
	//server url in host:port format
	Server string
	// storage to store data offline
	Storage       LocalStorage
	Client        *http.Client
	AuthCookie    *http.Cookie
	LoggedOffline bool
	Key           []byte
}

type LocalStorage interface {
	InitStorage(key []byte) error
	DeleteLocalStorage() error

	//Cards processing methods
	SaveCard(card database.Card, key []byte) error
	GetCard(cardname string, key []byte) (card database.Card, err error)
	ListCards() (cards []string, err error)

	//Login credentials processing methods
	SaveLoginCreds(logincreds database.LoginCreds, key []byte) error
	GetLoginCreds(logincredsname string, key []byte) (logincreds database.LoginCreds, err error)
	ListLoginCreds() (logincreds []string, err error)

	//Notes processing methods
	SaveNote(note database.Note, key []byte) error
	GetNote(notename string, key []byte) (note database.Note, err error)
	ListNotes() (notes []string, err error)

	//Binaries processing methods
	SaveBinary(binary database.Binary, key []byte) error
	GetBinary(binaryname string, key []byte) (binary database.Binary, err error)
	ListBinaries() (binaries []string, err error)
}

func NewClient() *Client {
	return &Client{
		Server:  config.ClientCfg.Address,
		Storage: localstorage.NewStorage(),
		Client:  &http.Client{},
	}
}

// PrintAvailableCommands function prints available commands
func PrintAvailableCommands(commands map[string]func([]string)) {
	fmt.Println("Available commands:")
	fmt.Println("- help")
	fmt.Println("- quit")
	for cmd := range commands {
		fmt.Println("-", cmd)
	}
}
