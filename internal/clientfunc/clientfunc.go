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
	SaveCard(card database.Card, key []byte) error
	GetCard(cardname string, key []byte) (card database.Card, err error)
	ListCards() (cards []string, err error)
	InitStorage(key []byte) error
	DeleteLocalStorage() error
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
