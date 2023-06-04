package clientfunc

import (
	"fmt"
	"net/http"

	"github.com/gambruh/gophkeeper/internal/config"
	"github.com/gambruh/gophkeeper/internal/database"
	"github.com/gambruh/gophkeeper/internal/localstorage"
)

type Client struct {
	Server        string
	Storage       LocalStorage
	Client        *http.Client
	AuthCookie    *http.Cookie
	LoggedOffline bool
	Key           []byte
}

type LocalStorage interface {
	SaveCard(card database.Card, key []byte) error
	GetCard(cardname string, key []byte) (card database.Card, err error)
	ListCards(key []byte) (cards []string, err error)
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
