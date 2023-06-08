package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/gambruh/gophkeeper/internal/clientfunc"
	"github.com/gambruh/gophkeeper/internal/compileinfo"
	"github.com/gambruh/gophkeeper/internal/config"
)

var (
	buildVersion string
	buildDate    string
	buildCommit  string
)

func main() {

	// вывод информации о компиляции
	compileinfo.PrintCompileInfo(buildVersion, buildDate, buildCommit)

	// Init client config
	config.InitClientFlags()
	config.SetClientConfig()

	//Init new client
	client := clientfunc.NewClient()

	// creating context for graceful shutdown
	ctxShutdown, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	// Ticker for synchronization
	syncTime := time.NewTicker(config.ClientCfg.CheckTime)
	defer syncTime.Stop()

	// Define available commands
	commands := map[string]func([]string){
		"register":  client.Register,
		"login":     client.Login,
		"setcard":   client.SetCardCommand,
		"getcard":   client.GetCardCommand,
		"listcards": client.ListCardsCommand,
	}

	// goroutine for data synchronization
	go client.DataChecker(ctxShutdown, syncTime)

	fmt.Println("write help to get commands list")
	reader := bufio.NewReader(os.Stdin)

	for {
		select {
		case <-ctxShutdown.Done():
			err := client.CheckCards()
			if err != nil {
				log.Println("error in CheckCards1:", err)
			}

			defer fmt.Println("Client exited!")
			return
		default:

		}
		fmt.Println("Enter a command:")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			continue
		}
		// Trim any leading/trailing whitespace and newline characters
		input = strings.TrimSpace(input)

		// Process the command
		if input == "quit" {
			fmt.Println("Exiting...")
			break
		}
		if input == "help" {
			clientfunc.PrintAvailableCommands(commands)
			continue
		}

		// Get the command name
		inpt := strings.Split(input, " ")
		if len(inpt) < 1 {
			fmt.Println("Please specify a command.")
			clientfunc.PrintAvailableCommands(commands)
		}

		command := inpt[0]

		if fn, ok := commands[command]; ok {
			fn(inpt)
		} else {
			fmt.Printf("Unknown command: %s\n", command)
			clientfunc.PrintAvailableCommands(commands)
		}
	}

	err := client.CheckCards()
	if err != nil {
		log.Println("error in CheckAll2 function returned from CheckCards:", err)
	}

	defer fmt.Println("Client exited!")
}
