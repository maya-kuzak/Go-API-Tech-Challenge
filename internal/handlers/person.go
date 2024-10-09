// all handlers for person func
package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/maya-kuzak/Go-API-Tech-Challenge/internal/models"
)

// Return all Person objects from the database.
func (h *RequestHandler) GetAllPeople(w http.ResponseWriter, r *http.Request) {
	var people []models.CompletePerson

	//get person data
	rows, err := h.DB.Query("SELECT * FROM person")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//insert person data into people slice
	for rows.Next() {
		var person models.CompletePerson
		err := rows.Scan(&person.ID, &person.FirstName, &person.LastName, &person.Type, &person.Age)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		//find courses for each person
		courseRows, err := h.DB.Query("SELECT course_id FROM person_course WHERE person_id = $1", person.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		for courseRows.Next() {
			var courseID uint
			err := courseRows.Scan(&courseID)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			person.Courses = append(person.Courses, courseID)
		}
		people = append(people, person)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(people); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
