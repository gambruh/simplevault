package main

import (
	"log"
	"net/http"

	"github.com/gambruh/gophkeeper/internal/auth"
	"github.com/gambruh/gophkeeper/internal/config"
	"github.com/gambruh/gophkeeper/internal/database"
	"github.com/gambruh/gophkeeper/internal/handlers"
)

func main() {
	config.InitFlags()
	config.SetConfig()
	authstorage := auth.GetAuthDB()
	defstorage := database.GetDB()

	service := handlers.NewService(defstorage, authstorage)

	server := &http.Server{
		Addr:    config.Cfg.Address,
		Handler: service.Service(),
	}

	log.Println(server.ListenAndServe())
}
