package main

import (
	"crypto/tls"
	"log"
	"net/http"
	"os"

	"github.com/gambruh/simplevault/internal/auth"
	"github.com/gambruh/simplevault/internal/config"
	"github.com/gambruh/simplevault/internal/handlers"
	"github.com/gambruh/simplevault/internal/storage/database"
)

func main() {

	// configuring logs
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	// setting file to write logs to
	file, err := os.OpenFile("server.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(file)

	// read CLI flags
	config.InitFlags()
	config.SetConfig()

	// connect to the authentication database
	authstorage := auth.GetAuthDB()

	// connect to the main storage
	defstorage := database.GetDB()

	// create a new service and
	service := handlers.NewService(defstorage, authstorage)

	server := &http.Server{
		Addr:      config.Cfg.Address,
		Handler:   service.Service(),
		TLSConfig: &tls.Config{},
	}

	log.Println(server.ListenAndServeTLS("cert.pem", "privatekey.pem"))
}
