// Package clientfunc contains all functions for the client to exchange data with the server, and to store it in localstorage
package clientfunc

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gambruh/gophkeeper/internal/config"
	"github.com/gambruh/gophkeeper/internal/database"
	"github.com/gambruh/gophkeeper/internal/localstorage"
)

type Client struct {
	//server url in host:port format
	Server string

	// storage interface to store data offline
	Storage LocalStorage

	// standard go-struct to make http-client
	Client *http.Client

	// this cookie will be applied to any http-request sent from the client, in case of successful online authentication
	AuthCookie *http.Cookie

	// box to be checked if user logged offline
	LoggedOffline bool

	// this is an encryption key
	Key []byte
}

// LocalStorage is an interfance
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

// NewClient function return new clientfunc.Client
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

// ResponseToCommand function parses an input and gets the command name.
// In case if command is valid it launches client methods to execute the command.
func (c *Client) ResponseToCommand(ctxShutdown context.Context, wgShutdown *sync.WaitGroup, quit chan<- struct{}, commands map[string]func([]string)) {
	defer wgShutdown.Done()
	for {
		select {
		case <-ctxShutdown.Done():
			return
		default:
			// going down
		}
		reader := bufio.NewReader(os.Stdin)

		fmt.Println()
		fmt.Println("Enter a command:")

		input, err := reader.ReadString('\n')
		if err != nil {
			return
		}
		// Trim any leading/trailing whitespace and newline characters
		input = strings.TrimSpace(input)

		// Process the command
		if input == "quit" {
			quit <- struct{}{}
			fmt.Println("Exiting...")
			return
		}
		if input == "help" {
			PrintAvailableCommands(commands)
			continue
		}

		// Get the command name
		inpt := strings.SplitN(input, " ", 2)
		if len(inpt) < 1 {
			fmt.Println("Please specify a command.")
			PrintAvailableCommands(commands)
		}
		command := inpt[0]

		// Execute the command
		if fn, ok := commands[command]; ok {
			fn(inpt)
		} else {
			fmt.Printf("Unknown command: %s\n", command)
			PrintAvailableCommands(commands)
		}
	}
}

// DataChecker synchronizes data between client and server
func (c *Client) DataChecker(context context.Context, wgShutdown *sync.WaitGroup, ticker *time.Ticker, quit <-chan struct{}) {
	defer wgShutdown.Done()
	for {
		select {
		case <-context.Done():
			return
		case <-quit:
			return
		case <-ticker.C:
			err := c.CheckAll()
			if err != nil {
				log.Println("error in DataChecker function returned from CheckAll:", err)
			}
		}
	}

}