package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/gambruh/simplevault/internal/clientfunc"
	"github.com/gambruh/simplevault/internal/config"
	"github.com/spf13/cobra"
)

func newRootCmd(run func(config.ClientConfig) error) *cobra.Command {
	flags := config.DefaultClientFlags()

	cmd := &cobra.Command{
		Use:          "client",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.SetClientConfig(flags)
			if err != nil {
				return err
			}
			return run(cfg)
		},
	}

	cmd.Flags().StringVarP(&flags.Address, "address", "a", flags.Address, "server address in format host:port")
	cmd.Flags().StringVarP(&flags.ClientCert, "clientcert", "s", flags.ClientCert, "path to client's certificate file")
	cmd.Flags().StringVarP(&flags.PrivateKey, "privatekey", "p", flags.PrivateKey, "path to file with public key for agent")
	cmd.Flags().StringVar(&flags.LocalStorage, "localstorage", flags.LocalStorage, "address of the folder to store files")
	cmd.Flags().DurationVarP(&flags.CheckTime, "checktime", "t", flags.CheckTime, "interval in time.Duration format (10s, 5m) to check data from DB")
	cmd.Flags().StringVar(&flags.BinInputFolder, "bininputfolder", flags.BinInputFolder, "folder to put binaries in to be sent")
	cmd.Flags().StringVar(&flags.BinOutputFolder, "binoutputfolder", flags.BinOutputFolder, "folder to store received binaries")

	return cmd
}

func runClient(cfg config.ClientConfig) error {
	client := clientfunc.NewClient(cfg)

	ctxShutdown, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	quit := make(chan struct{})

	var wgShutdown sync.WaitGroup
	wgShutdown.Add(2)

	syncTime := time.NewTicker(config.ClientCfg.CheckTime)
	defer syncTime.Stop()

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

	go client.DataChecker(ctxShutdown, &wgShutdown, syncTime, quit)

	fmt.Println("Write help to get commands list")
	go client.ResponseToCommand(ctxShutdown, &wgShutdown, quit, commands)

	wgShutdown.Wait()
	if err := client.CheckAll(); err != nil {
		log.Println("error in CheckAll function:", err)
	}

	fmt.Println("Client exited!")
	return nil
}
