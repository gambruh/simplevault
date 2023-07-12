package main

import (
	"crypto/tls"
	"log"
	"net/http"

	"github.com/gambruh/gophkeeper/internal/auth"
	"github.com/gambruh/gophkeeper/internal/config"
	"github.com/gambruh/gophkeeper/internal/handlers"
	"github.com/gambruh/gophkeeper/internal/storage/database"
)

func main() {
	config.InitFlags()
	config.SetConfig()
	authstorage := auth.GetAuthDB()
	defstorage := database.GetDB()

	service := handlers.NewService(defstorage, authstorage)

	server := &http.Server{
		Addr:      config.Cfg.Address,
		Handler:   service.Service(),
		TLSConfig: &tls.Config{},
	}

	log.Println(server.ListenAndServeTLS("cert.pem", "privatekey.pem"))
}
