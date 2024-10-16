//crete webserver using chi (new router, etc)

package webserver

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/maya-kuzak/Go-API-Tech-Challenge/internal/handlers"
	"github.com/maya-kuzak/Go-API-Tech-Challenge/internal/routes"
)

func NewServer(db *sql.DB) {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	handler := &handlers.RequestHandler{DB: db}
	routes.GetRoutes(r, handler)

	log.Print("\nStarting server on port :8000\n")
	err := http.ListenAndServe("localhost:8000", r)
	if err != nil {
		log.Fatal("Listen and server error: ", err)
	}

}
