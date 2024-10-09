package main

import (
	"log"

	"github.com/maya-kuzak/Go-API-Tech-Challenge/internal/database"
	"github.com/maya-kuzak/Go-API-Tech-Challenge/internal/webserver"
)

func main() {
	db, err := database.Init()
	if err != nil {
		log.Fatal(err)
	}

	webserver.NewServer(db)

}
