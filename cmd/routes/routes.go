// all api calls
package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/maya-kuzak/Go-API-Tech-Challenge/cmd/handlers"
)

func GetRoutes(r chi.Router, handler *handlers.RequestHandler) {
	// course routes
	r.Get("/api/course", handler.GetAllCourses)
	r.Get("/api/course/{id}", handler.GetCourse)
	r.Put("/api/course/{id}", handler.UpdateCourse)
	r.Post("/api/course", handler.CreateCourse)
	r.Delete("/api/course/{id}", handler.DeleteCourse)

	// person routes
	r.Get("/api/person", handler.GetAllPeople)     //takes querys of name (first or last) and age
	r.Get("/api/person/{name}", handler.GetPerson) // name = first + ' ' + last
	r.Put("/api/person/{name}", handler.UpdatePerson)
	r.Post("/api/person", handler.CreatePerson)
	r.Delete("/api/person/{name}", handler.DeletePerson)

}
