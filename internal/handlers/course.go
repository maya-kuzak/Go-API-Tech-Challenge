// all handlers for course func
package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/maya-kuzak/Go-API-Tech-Challenge/internal/models"
)

type RequestHandler struct {
	DB *sql.DB
}

// Return all Course objects from the database.
func (h *RequestHandler) GetAllCourses(w http.ResponseWriter, r *http.Request) {
	var courses []models.Course

	rows, err := h.DB.Query("SELECT * FROM course")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for rows.Next() {
		var course models.Course
		err := rows.Scan(&course.ID, &course.Name)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		courses = append(courses, course)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(courses); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func (h *RequestHandler) GetCourse(w http.ResponseWriter, r *http.Request) {
	var course models.Course
	id := chi.URLParam(r, "id")

	intID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	row := h.DB.QueryRow("SELECT * FROM course WHERE id = $1", intID)
	err = row.Scan(&course.ID, &course.Name)
	if err != nil {
		http.Error(w, "Invalid course ID", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(course); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *RequestHandler) UpdateCourse(w http.ResponseWriter, r *http.Request) {
	var course models.Course
	id := chi.URLParam(r, "id")

	intID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid Course ID", http.StatusInternalServerError)
		return
	}

	err = json.NewDecoder(r.Body).Decode(&course)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusInternalServerError)
		return
	}

	// Validate the Course object
	if course.Name == "" {
		http.Error(w, "Course name is required", http.StatusBadRequest)
		return
	}

	//update course
	_, err = h.DB.Exec("UPDATE course SET name = $1 WHERE id = $2", course.Name, intID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//return updated course
	row := h.DB.QueryRow("SELECT * FROM course WHERE id = $1", intID)
	err = row.Scan(&course.ID, &course.Name)
	if err != nil {
		http.Error(w, "Invalid course ID", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(course); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *RequestHandler) CreateCourse(w http.ResponseWriter, r *http.Request) {
	var course models.Course

	err := json.NewDecoder(r.Body).Decode(&course)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusInternalServerError)
		return
	}

	// Validate the Course object
	if course.Name == "" {
		http.Error(w, "Course name is required", http.StatusBadRequest)
		return
	}

	err = h.DB.QueryRow("INSERT INTO course (name) VALUES ($1) RETURNING id", course.Name).Scan(&course.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(course); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *RequestHandler) DeleteCourse(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	intID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid Course ID", http.StatusInternalServerError)
		return
	}

	//delete course
	_, err = h.DB.Exec("DELETE FROM course WHERE id = $1", intID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)

	// Return a JSON-formatted deletion confirmation message
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]string{"message": "Course deleted successfully"}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Error encoding response: "+err.Error(), http.StatusInternalServerError)
		return
	}
}
