package main

import (
	"flag"
	"fmt"

	"github.com/gambruh/gophkeeper/internal/clientfunc"
	"github.com/gambruh/gophkeeper/internal/config"
)

func main() {

	// Init config
	config.InitClientFlags()
	config.SetClientConfig()
	fmt.Printf("%v\n", config.ClientCfg)

	//Init new client
	client := clientfunc.NewClient()

	// Define available commands
	commands := map[string]func([]string){
		"register": client.Register,
		"login":    client.Login,
		"setcard":  client.SetCardCommand,
	}

	// Parse command-line arguments
	flag.Parse()

	// Get the command name
	args := flag.Args()
	if len(args) < 1 {
		fmt.Println("Please specify a command.")
		clientfunc.PrintAvailableCommands(commands)
		return
	}

	command := args[0]

	if fn, ok := commands[command]; ok {
		fn(args)
	} else {
		fmt.Printf("Unknown command: %s\n", command)
		clientfunc.PrintAvailableCommands(commands)
		return
	}
}
