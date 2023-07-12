package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
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

	// Set client config
	cfg := config.SetClientConfig()

	//Init new client
	client := clientfunc.NewClient(cfg)

	// creating context for graceful shutdown
	ctxShutdown, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	// channel for normal exit from the user
	quit := make(chan struct{})

	var wgShutdown sync.WaitGroup
	wgShutdown.Add(2)

	// Ticker for synchronization
	syncTime := time.NewTicker(config.ClientCfg.CheckTime)
	defer syncTime.Stop()

	// Define available commands
	commands := map[string]func([]string){
		"register":       client.Register,
		"login":          client.Login,
		"setcard":        client.SetCardCommand,
		"getcard":        client.GetCardCommand,
		"listcards":      client.ListCardsCommand,
		"setlogincreds":  client.SetLoginCredsCommand,
		"getlogincreds":  client.GetLoginCredsCommand,
		"listlogincreds": client.ListLoginCredsCommand,
		"setnote":        client.SetNoteCommand,
		"getnote":        client.GetNoteCommand,
		"listnotes":      client.ListNotesCommand,
		"setbinary":      client.SetBinaryCommand,
		"getbinary":      client.GetBinaryCommand,
		"listbinaries":   client.ListBinariesCommand,
	}

	// goroutine for data synchronization between client and server
	go client.DataChecker(ctxShutdown, &wgShutdown, syncTime, quit)

	// goroutine for command recognition and client responding with actions
	fmt.Println("Write help to get commands list")
	go client.ResponseToCommand(ctxShutdown, &wgShutdown, quit, commands)

	wgShutdown.Wait()
	err := client.CheckAll()
	if err != nil {
		log.Println("error in CheckAll function:", err)
	}

	defer fmt.Println("Client exited!")
}
