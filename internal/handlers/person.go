// all handlers for person func
package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/maya-kuzak/Go-API-Tech-Challenge/internal/models"
)

// Return all Person objects from the database.
func (h *RequestHandler) GetAllPeople(w http.ResponseWriter, r *http.Request) {
	var people []models.CompletePerson

	//get query params
	name := r.URL.Query().Get("name")
	age := r.URL.Query().Get("age")

	query := "SELECT id, first_name, last_name, type, age FROM person"
	var args []interface{}
	var conditions []string

	if name != "" {
		conditions = append(conditions, "first_name = $1 OR last_name = $1")
		args = append(args, name)
	}

	if age != "" {
		placeholder := "$1"
		if name != "" {
			placeholder = "$2"
		}
		conditions = append(conditions, "age = "+placeholder)
		args = append(args, age)
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	//get person data
	rows, err := h.DB.Query(query, args...)
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

// Return a given Person from the database.
func (h *RequestHandler) GetPerson(w http.ResponseWriter, r *http.Request) {
	var person models.CompletePerson

	//get query params
	fullName := chi.URLParam(r, "name")

	query := "SELECT id, first_name, last_name, type, age FROM person WHERE first_name || ' ' || last_name = $1"
	args := []interface{}{fullName}

	// Get person data
	row := h.DB.QueryRow(query, args...)
	err := row.Scan(&person.ID, &person.FirstName, &person.LastName, &person.Type, &person.Age)
	if err != nil {
		http.Error(w, "Person not found"+err.Error(), http.StatusInternalServerError)
		return
	}

	//find courses for each person
	courseRows, err := h.DB.Query("SELECT course_id FROM person_course WHERE person_id = $1", person.ID)
	if err != nil {
		http.Error(w, "Error fetching courses: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer courseRows.Close()

	for courseRows.Next() {
		var courseID uint
		err := courseRows.Scan(&courseID)
		if err != nil {
			http.Error(w, "Error scanning course ID"+err.Error(), http.StatusInternalServerError)
			return
		}
		person.Courses = append(person.Courses, courseID)
	}

	if err := courseRows.Err(); err != nil {
		http.Error(w, "Error iterating over courses: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(person); err != nil {
		http.Error(w, "Error encoding response"+err.Error(), http.StatusInternalServerError)
		return
	}
}

// Update a given Person in the database based on name.
func (h *RequestHandler) UpdatePerson(w http.ResponseWriter, r *http.Request) {
	var updatedPerson models.CompletePerson

	// Parse the JSON request body
	if err := json.NewDecoder(r.Body).Decode(&updatedPerson); err != nil {
		http.Error(w, "Invalid request payload: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Validate the Person object
	if updatedPerson.FirstName == "" || updatedPerson.LastName == "" || updatedPerson.Type == "" || updatedPerson.Age == 0 {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	// Get path param
	fullName := chi.URLParam(r, "name")

	// Construct the SQL UPDATE query
	query := `
        UPDATE person
        SET first_name = $1, last_name = $2, type = $3, age = $4
        WHERE first_name || ' ' || last_name = $5
        RETURNING id, first_name, last_name, type, age
    `
	args := []interface{}{updatedPerson.FirstName, updatedPerson.LastName, updatedPerson.Type, updatedPerson.Age, fullName}

	// Execute the UPDATE query and fetch the updated person
	row := h.DB.QueryRow(query, args...)
	err := row.Scan(&updatedPerson.ID, &updatedPerson.FirstName, &updatedPerson.LastName, &updatedPerson.Type, &updatedPerson.Age)
	if err != nil {
		http.Error(w, "Error updating person: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Update courses for the person
	_, err = h.DB.Exec("DELETE FROM person_course WHERE person_id = $1", updatedPerson.ID)
	if err != nil {
		http.Error(w, "Error deleting old courses: "+err.Error(), http.StatusInternalServerError)
		return
	}

	for _, courseID := range updatedPerson.Courses {
		_, err := h.DB.Exec("INSERT INTO person_course (person_id, course_id) VALUES ($1, $2)", updatedPerson.ID, courseID)
		if err != nil {
			http.Error(w, "Error inserting new courses: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// Fetch the updated courses
	courseRows, err := h.DB.Query("SELECT course_id FROM person_course WHERE person_id = $1", updatedPerson.ID)
	if err != nil {
		http.Error(w, "Error fetching courses: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer courseRows.Close()

	updatedPerson.Courses = nil // Reset courses slice
	for courseRows.Next() {
		var courseID uint
		err := courseRows.Scan(&courseID)
		if err != nil {
			http.Error(w, "Error scanning course ID: "+err.Error(), http.StatusInternalServerError)
			return
		}
		updatedPerson.Courses = append(updatedPerson.Courses, courseID)
	}

	if err := courseRows.Err(); err != nil {
		http.Error(w, "Error iterating over courses: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the updated Person object as a JSON response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(updatedPerson); err != nil {
		http.Error(w, "Error encoding response: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

// Create a new Person in the database.
func (h *RequestHandler) CreatePerson(w http.ResponseWriter, r *http.Request) {
	var newPerson models.CompletePerson

	// Parse the JSON request body
	if err := json.NewDecoder(r.Body).Decode(&newPerson); err != nil {
		http.Error(w, "Invalid request payload: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Validate the Person object
	if newPerson.FirstName == "" || newPerson.LastName == "" || newPerson.Type == "" || newPerson.Age == 0 {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	// Construct the SQL INSERT query
	query := `
        INSERT INTO person (first_name, last_name, type, age)
        VALUES ($1, $2, $3, $4)
        RETURNING id
    `
	args := []interface{}{newPerson.FirstName, newPerson.LastName, newPerson.Type, newPerson.Age}

	// Execute the INSERT query and fetch the new person's ID
	var newPersonID uint
	err := h.DB.QueryRow(query, args...).Scan(&newPersonID)
	if err != nil {
		http.Error(w, "Error creating person: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Update courses for the new person
	for _, courseID := range newPerson.Courses {
		// Check if the course_id exists in the course table
		var exists bool
		err := h.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM course WHERE id = $1)", courseID).Scan(&exists)
		if err != nil {
			http.Error(w, "Error checking course existence: "+err.Error(), http.StatusInternalServerError)
			return
		}
		if !exists {
			http.Error(w, "Course ID does not exist: "+string(courseID), http.StatusBadRequest)
			return
		}

		_, err = h.DB.Exec("INSERT INTO person_course (person_id, course_id) VALUES ($1, $2)", newPersonID, courseID)
		if err != nil {
			http.Error(w, "Error inserting new courses: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// Return the new Person object's ID as a JSON response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]uint{"id": newPersonID}); err != nil {
		http.Error(w, "Error encoding response: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

// Delete a given Person from the database based on name.
func (h *RequestHandler) DeletePerson(w http.ResponseWriter, r *http.Request) {
	// Get path param
	fullName := chi.URLParam(r, "name")

	// Find the person ID based on the full name
	var personID int
	err := h.DB.QueryRow("SELECT id FROM person WHERE first_name || ' ' || last_name = $1", fullName).Scan(&personID)
	if err != nil {
		http.Error(w, "Error finding person: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Delete associated records in the person_course table
	_, err = h.DB.Exec("DELETE FROM person_course WHERE person_id = $1", personID)
	if err != nil {
		http.Error(w, "Error deleting associated courses: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Delete the person
	_, err = h.DB.Exec("DELETE FROM person WHERE id = $1", personID)
	if err != nil {
		http.Error(w, "Error deleting person: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Return a success message as a JSON response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]string{"message": "Person deleted successfully"}); err != nil {
		http.Error(w, "Error encoding response: "+err.Error(), http.StatusInternalServerError)
		return
	}
}
