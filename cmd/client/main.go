package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/gambruh/gophkeeper/internal/clientfunc"
	"github.com/gambruh/gophkeeper/internal/config"
)

func main() {
	// Init config
	config.InitClientFlags()
	config.SetClientConfig()

	//Init new client
	client := clientfunc.NewClient()

	// Define available commands
	commands := map[string]func([]string){
		"register": client.Register,
		"login":    client.Login,
		"setcard":  client.SetCardCommand,
		"getcard":  client.GetCardCommand,
	}

	fmt.Println("write help to get commands list")
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("Enter a command: ")
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
	defer fmt.Println("Client exited!")
}
