package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

// TestGetAllCourses tests the GetAllCourses handler.
func TestGetAllCourses(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	handler := &RequestHandler{DB: db}

	// Mock the database response
	rows := sqlmock.NewRows([]string{"id", "name"}).
		AddRow(1, "Course 1").
		AddRow(2, "Course 2")
	mock.ExpectQuery("SELECT \\* FROM course").WillReturnRows(rows)

	// Create a new HTTP request
	req, err := http.NewRequest("GET", "/courses", nil)
	assert.NoError(t, err)

	// Record the HTTP response
	rr := httptest.NewRecorder()
	handler.GetAllCourses(rr, req)

	// Assert the response
	assert.Equal(t, http.StatusOK, rr.Code)
	var courses []Course
	err = json.NewDecoder(rr.Body).Decode(&courses)
	assert.NoError(t, err)
	assert.Len(t, courses, 2)
	assert.Equal(t, "Course 1", courses[0].Name)
	assert.Equal(t, "Course 2", courses[1].Name)
}

// TestGetCourse tests the GetCourse handler.
func TestGetCourse(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	handler := &RequestHandler{DB: db}

	// Mock the database response
	row := sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "Course 1")
	mock.ExpectQuery("SELECT \\* FROM course WHERE id = \\$1").WithArgs(1).WillReturnRows(row)

	// Create a new HTTP request
	req, err := http.NewRequest("GET", "/courses/1", nil)
	assert.NoError(t, err)

	// Set the URL parameter
	rr := httptest.NewRecorder()
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	// Call the handler
	handler.GetCourse(rr, req)

	// Assert the response
	assert.Equal(t, http.StatusOK, rr.Code)
	var course Course
	err = json.NewDecoder(rr.Body).Decode(&course)
	assert.NoError(t, err)
	assert.Equal(t, "Course 1", course.Name)
}

// TestUpdateCourse tests the UpdateCourse handler.
func TestUpdateCourse(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	handler := &RequestHandler{DB: db}

	// Create a course object and marshal it to JSON
	course := Course{Name: "Updated Course"}
	courseJSON, _ := json.Marshal(course)

	// Mock the database response
	mock.ExpectExec("UPDATE course SET name = \\$1 WHERE id = \\$2").
		WithArgs(course.Name, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	row := sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "Updated Course")
	mock.ExpectQuery("SELECT \\* FROM course WHERE id = \\$1").WithArgs(1).WillReturnRows(row)

	// Create a new HTTP request
	req, err := http.NewRequest("PUT", "/courses/1", bytes.NewBuffer(courseJSON))
	assert.NoError(t, err)

	// Set the URL parameter
	rr := httptest.NewRecorder()
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	// Call the handler
	handler.UpdateCourse(rr, req)

	// Assert the response
	assert.Equal(t, http.StatusOK, rr.Code)
	var updatedCourse Course
	err = json.NewDecoder(rr.Body).Decode(&updatedCourse)
	assert.NoError(t, err)
	assert.Equal(t, "Updated Course", updatedCourse.Name)
}

// TestCreateCourse tests the CreateCourse handler.
func TestCreateCourse(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	handler := &RequestHandler{DB: db}

	// Create a course object and marshal it to JSON
	course := Course{Name: "New Course"}
	courseJSON, _ := json.Marshal(course)

	// Mock the database response
	mock.ExpectQuery("INSERT INTO course \\(name\\) VALUES \\(\\$1\\) RETURNING id").
		WithArgs(course.Name).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	// Create a new HTTP request
	req, err := http.NewRequest("POST", "/courses", bytes.NewBuffer(courseJSON))
	assert.NoError(t, err)

	// Record the HTTP response
	rr := httptest.NewRecorder()
	handler.CreateCourse(rr, req)

	// Assert the response
	assert.Equal(t, http.StatusOK, rr.Code)
	var newCourse Course
	err = json.NewDecoder(rr.Body).Decode(&newCourse)
	assert.NoError(t, err)
	assert.Equal(t, "New Course", newCourse.Name)
}

// TestDeleteCourse tests the DeleteCourse handler.
func TestDeleteCourse(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	handler := &RequestHandler{DB: db}

	// Mock the database response
	mock.ExpectExec("DELETE FROM course WHERE id = \\$1").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Create a new HTTP request
	req, err := http.NewRequest("DELETE", "/courses/1", nil)
	assert.NoError(t, err)

	// Set the URL parameter
	rr := httptest.NewRecorder()
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	// Call the handler
	handler.DeleteCourse(rr, req)

	// Assert the response
	assert.Equal(t, http.StatusNoContent, rr.Code)
}
