package main

import (
	"crypto/tls"
	"log"
	"net/http"

	"github.com/gambruh/simplevault/internal/auth"
	"github.com/gambruh/simplevault/internal/config"
	"github.com/gambruh/simplevault/internal/handlers"
	"github.com/gambruh/simplevault/internal/storage/database"
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
