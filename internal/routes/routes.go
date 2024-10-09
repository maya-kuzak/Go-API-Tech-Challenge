// all api calls
package routes

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/maya-kuzak/Go-API-Tech-Challenge/internal/handlers"
)

func GetRoutes(r chi.Router, handler *handlers.RequestHandler) {
	// course routes
	r.Get("/api/course", handler.GetAllCourses)
	r.Get("/api/course/{id}", handler.GetCourse)
	r.Put("/api/course/{id}", handler.UpdateCourse)
	r.Post("/api/course", handler.CreateCourse)
	r.Delete("/api/course/{id}", handler.DeleteCourse)

	// person routes
	r.Get("/api/person", handler.GetAllPeople)
	r.Get("/api/person/{name}", getPerson)
	r.Put("/api/person/{name}", updatePerson)
	r.Post("/api/person", createPerson)
	r.Delete("/api/person/{name}", deletePerson)

}

// Placeholder handler functions

func getPerson(w http.ResponseWriter, r *http.Request)    {}
func updatePerson(w http.ResponseWriter, r *http.Request) {}
func createPerson(w http.ResponseWriter, r *http.Request) {}
func deletePerson(w http.ResponseWriter, r *http.Request) {}
