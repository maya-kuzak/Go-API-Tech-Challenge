//crete webserver using chi (new router, etc)

package webserver

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func NewServer() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.ListenAndServe("localhost:8000", r)
}
