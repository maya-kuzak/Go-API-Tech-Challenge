//crete webserver using chi (new router, etc)

package webserver

import (
	"log"
	"net/http"

	"gorm.io/gorm"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewServer(db *gorm.DB) {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	log.Print("\nStarting server on port :8000\n")
	err := http.ListenAndServe("localhost:8000", r)
	if err != nil {
		log.Fatal("Listen and server error: ", err)
	}

}
