package routes

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-chi/chi/v5"
	"github.com/maya-kuzak/Go-API-Tech-Challenge/internal/handlers"
	"github.com/stretchr/testify/assert"
)

func TestGetRoutes(t *testing.T) {

	// Create a new mock database connection
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Create a new request handler
	handler := &handlers.RequestHandler{DB: db}
	r := chi.NewRouter()
	GetRoutes(r, handler)

	// Define the test cases with HTTP methods, URLs, and expected status codes
	tests := []struct {
		method       string
		url          string
		expectedCode int
	}{
		{"GET", "/api/course", http.StatusOK},
		{"GET", "/api/course/1", http.StatusOK},
		{"GET", "/api/person", http.StatusOK},
		{"GET", "/api/person/John Doe", http.StatusOK},

		{"PUT", "/api/course/1", http.StatusOK},
		{"POST", "/api/course", http.StatusOK},
		{"DELETE", "/api/course/1", http.StatusNoContent},
	}

	// Mock the database responses
	mock.ExpectQuery("SELECT \\* FROM course").WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "Course 1").AddRow(2, "Course 2"))
	mock.ExpectQuery("SELECT \\* FROM course WHERE id = \\$1").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "Course 1"))
	mock.ExpectQuery("SELECT id, first_name, last_name, type, age FROM person").WillReturnRows(sqlmock.NewRows([]string{"id", "first_name", "last_name", "type", "age"}).AddRow(1, "John", "Doe", "student", 25).AddRow(2, "Jane", "Doe", "professor", 30))

	mock.ExpectQuery("SELECT course_id FROM person_course WHERE person_id = \\$1").
		WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"course_id"}).AddRow(1).AddRow(2))
	mock.ExpectQuery("SELECT course_id FROM person_course WHERE person_id = \\$1").
		WithArgs(2).WillReturnRows(sqlmock.NewRows([]string{"course_id"}).AddRow(1).AddRow(3))
	mock.ExpectQuery("SELECT id, first_name, last_name, type, age FROM person WHERE first_name \\|\\| ' ' \\|\\| last_name = \\$1").WithArgs("John Doe").WillReturnRows(sqlmock.NewRows([]string{"id", "first_name", "last_name", "type", "age"}).AddRow(1, "John", "Doe", "student", 25))
	mock.ExpectQuery("SELECT course_id FROM person_course WHERE person_id = \\$1").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"course_id"}).AddRow(1))

	mock.ExpectExec("UPDATE course SET name = \\$1 WHERE id = \\$2").WithArgs("Updated Course", 1).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectQuery("SELECT \\* FROM course WHERE id = \\$1").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "Updated Course"))

	mock.ExpectQuery("INSERT INTO course \\(name\\) VALUES \\(\\$1\\) RETURNING id").
		WithArgs("New Course").WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	mock.ExpectExec("DELETE FROM course WHERE id = \\$1").
		WithArgs(1).WillReturnResult(sqlmock.NewResult(1, 1))

	// Run the test cases
	for _, tt := range tests {
		t.Run(tt.method+"_"+tt.url, func(t *testing.T) {
			var req *http.Request
			var err error

			// Create a new HTTP request
			if tt.method == "PUT" && tt.url == "/api/course/1" {
				// Create a course object and marshal it to JSON
				course := handlers.Course{Name: "Updated Course"}
				courseJSON, _ := json.Marshal(course)
				req, err = http.NewRequest(tt.method, tt.url, bytes.NewBuffer(courseJSON))
				req.Header.Set("Content-Type", "application/json")
			} else if tt.method == "POST" && tt.url == "/api/course" {
				// Create a course object and marshal it to JSON
				course := handlers.Course{Name: "New Course"}
				courseJSON, _ := json.Marshal(course)
				req, err = http.NewRequest(tt.method, tt.url, bytes.NewBuffer(courseJSON))
				req.Header.Set("Content-Type", "application/json")
			} else if tt.method == "DELETE" && tt.url == "/api/course/1" {
				req, err = http.NewRequest(tt.method, tt.url, nil)
			} else {
				req, err = http.NewRequest(tt.method, tt.url, nil)
			}

			// Assert there were no errors creating the request
			assert.NoError(t, err)

			// Create a new recorder to capture response
			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)

			// Assert the response status code matches expected code
			assert.Equal(t, tt.expectedCode, rr.Code)
		})
	}

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
