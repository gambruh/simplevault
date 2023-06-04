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
}

type LocalStorage interface {
	SaveCard(card database.Card) error
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
