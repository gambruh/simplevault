package main

import (
	"crypto/tls"
	"net/http"

	"github.com/gambruh/simplevault/internal/auth"
	"github.com/gambruh/simplevault/internal/config"
	"github.com/gambruh/simplevault/internal/handlers"
	"github.com/gambruh/simplevault/internal/storage/database"
	"github.com/spf13/cobra"
)

func newRootCmd(run func() error) *cobra.Command {
	flags := config.DefaultServerFlags()

	cmd := &cobra.Command{
		Use:          "server",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := config.SetConfig(flags); err != nil {
				return err
			}
			return run()
		},
	}

	cmd.Flags().StringVarP(&flags.Address, "address", "a", flags.Address, "server address in format host:port")
	cmd.Flags().StringVar(&flags.Certificate, "cert", flags.Certificate, "certificate to run TLS")
	cmd.Flags().StringVar(&flags.PrivateKey, "privatekey", flags.PrivateKey, "server's private key")
	cmd.Flags().StringVarP(&flags.Database, "database", "d", flags.Database, "postgres database uri")
	cmd.Flags().StringVarP(&flags.Key, "key", "k", flags.Key, "key to hash")

	return cmd
}

func runServer() error {
	authstorage := auth.GetAuthDB()
	defstorage := database.GetDB()

	service := handlers.NewService(defstorage, authstorage)

	server := &http.Server{
		Addr:      config.Cfg.Address,
		Handler:   service.Service(),
		TLSConfig: &tls.Config{},
	}

	return server.ListenAndServeTLS(config.Cfg.Certificate, config.Cfg.PrivateKey)
}
