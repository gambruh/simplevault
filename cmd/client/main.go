package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"
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
	var wgShutdown sync.WaitGroup
	wgShutdown.Add(2)

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
	go client.DataChecker(ctxShutdown, syncTime, &wgShutdown)

	fmt.Println("write help to get commands list")

	go func() {
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
				fmt.Println("Exiting...")
				return
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
	}()

	wgShutdown.Wait()
	err := client.CheckCards()
	if err != nil {
		log.Println("error in CheckAll2 function returned from CheckCards:", err)
	}

	defer fmt.Println("Client exited!")
}
