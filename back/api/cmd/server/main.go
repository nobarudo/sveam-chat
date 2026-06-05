package main

import (
	"log"

	"svem-chat-api/internal/platform/config"
	"svem-chat-api/internal/platform/database"
	"svem-chat-api/internal/platform/server"
)

func main() {
	conf := config.GetConfig()

	db, err := database.Connect(&conf)
	if err != nil {
		log.Fatal(err)
	}
	database.Migration(db)
	server.Start(db)
}
