package main

import (
	"log"

	"github.com/maya-kuzak/Go-API-Tech-Challenge/internal/database"
	"github.com/maya-kuzak/Go-API-Tech-Challenge/internal/webserver"
)

//import database and webserver

func main() {
	//init databse
	db, err := database.Init()
	if err != nil {
		log.Fatal(err)
	}

	//pass database connection to new webserver
	webserver.NewServer(db)
}
