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

func TestGetAllPeople(t *testing.T) {
	// Create a new mock database connection
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Create a new request handler
	handler := &RequestHandler{DB: db}

	// Mock the database response for people
	// . - method chaining to be executed at same time (ex sqlmock.NewRows().AddRow())
	rows := sqlmock.NewRows([]string{"id", "first_name", "last_name", "type", "age"}).
		AddRow(1, "John", "Doe", "student", 25).
		AddRow(2, "Jane", "Smith", "professor", 30)
	mock.ExpectQuery("SELECT id, first_name, last_name, type, age FROM person").WillReturnRows(rows)

	// Mock the database response for courses for each person
	courseRows1 := sqlmock.NewRows([]string{"course_id"}).
		AddRow(1).
		AddRow(2)
	mock.ExpectQuery("SELECT course_id FROM person_course WHERE person_id = \\$1").WithArgs(1).WillReturnRows(courseRows1)

	courseRows2 := sqlmock.NewRows([]string{"course_id"}).
		AddRow(3)
	mock.ExpectQuery("SELECT course_id FROM person_course WHERE person_id = \\$1").WithArgs(2).WillReturnRows(courseRows2)

	// Create a new HTTP request
	req, err := http.NewRequest("GET", "/people", nil)
	assert.NoError(t, err)

	// Record the HTTP response
	rr := httptest.NewRecorder()
	handler.GetAllPeople(rr, req)

	// Assert the response status code
	assert.Equal(t, http.StatusOK, rr.Code)

	// Decode the response body into a slice of CompletePerson
	var people []CompletePerson
	err = json.NewDecoder(rr.Body).Decode(&people)
	assert.NoError(t, err)

	// Assert the number of people and their details
	//Assume if first name matches, rest of person data matches
	assert.Len(t, people, 2)
	assert.Equal(t, "John", people[0].FirstName)
	assert.Equal(t, "Jane", people[1].FirstName)

	// Assert the courses for each person
	assert.Len(t, people[0].Courses, 2)
	assert.Equal(t, uint(1), people[0].Courses[0])
	assert.Equal(t, uint(2), people[0].Courses[1])

	assert.Len(t, people[1].Courses, 1)
	assert.Equal(t, uint(3), people[1].Courses[0])
}

func TestGetPerson(t *testing.T) {
	//very similar to TestGetAllPeople

	//create new mock database connection
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	//create new request handler
	handler := &RequestHandler{DB: db}

	// Mock the database response
	row := sqlmock.NewRows([]string{"id", "first_name", "last_name", "type", "age"}).
		AddRow(1, "John", "Doe", "student", 25)
	mock.ExpectQuery("SELECT id, first_name, last_name, type, age FROM person WHERE first_name \\|\\| ' ' \\|\\| last_name = \\$1").
		WithArgs("John Doe").WillReturnRows(row)

	// Mock the database response for courses
	courseRows := sqlmock.NewRows([]string{"course_id"}).
		AddRow(1)
	mock.ExpectQuery("SELECT course_id FROM person_course WHERE person_id = \\$1").WithArgs(1).WillReturnRows(courseRows)

	req, err := http.NewRequest("GET", "/person/John Doe", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("name", "John Doe")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	handler.GetPerson(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	var person CompletePerson
	err = json.NewDecoder(rr.Body).Decode(&person)
	assert.NoError(t, err)
	assert.Equal(t, "John", person.FirstName)
}
func TestUpdatePerson(t *testing.T) {
	//mock database connection
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	handler := &RequestHandler{DB: db}

	// Mock the database response for finding the person ID
	row := sqlmock.NewRows([]string{"id"}).
		AddRow(1)
	mock.ExpectQuery("SELECT id FROM person WHERE first_name \\|\\| ' ' \\|\\| last_name = \\$1").
		WithArgs("John Doe").WillReturnRows(row)

	// Mock the database response for updating the person
	mock.ExpectExec("UPDATE person SET first_name = \\$1, last_name = \\$2, type = \\$3, age = \\$4 WHERE id = \\$5").
		WithArgs("John", "Doe", "student", 25, 1).WillReturnResult(sqlmock.NewResult(1, 1))

	// Mock the database response for deleting old courses
	mock.ExpectExec("DELETE FROM person_course WHERE person_id = \\$1").
		WithArgs(1).WillReturnResult(sqlmock.NewResult(1, 1))

	// Mock the database response for checking course existence
	mock.ExpectQuery("SELECT EXISTS\\(SELECT 1 FROM course WHERE id = \\$1\\)").
		WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

	// Mock the database response for inserting new courses
	mock.ExpectExec("INSERT INTO person_course \\(person_id, course_id\\) VALUES \\(\\$1, \\$2\\)").
		WithArgs(1, 1).WillReturnResult(sqlmock.NewResult(1, 1))

	person := CompletePerson{
		FirstName: "John",
		LastName:  "Doe",
		Type:      "student",
		Age:       25,
		Courses:   []uint{1},
	}
	body, err := json.Marshal(person)
	assert.NoError(t, err)

	req, err := http.NewRequest("PUT", "/person/John Doe", bytes.NewBuffer(body))
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("name", "John Doe")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	handler.UpdatePerson(rr, req)

	// Assert the response
	assert.Equal(t, http.StatusOK, rr.Code)
	var updatedPerson CompletePerson
	err = json.NewDecoder(rr.Body).Decode(&updatedPerson)
	assert.NoError(t, err)
	assert.Equal(t, "John", updatedPerson.FirstName)
}

func TestCreatePerson(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	handler := &RequestHandler{DB: db}

	// Mock the database response
	mock.ExpectQuery("INSERT INTO person \\(first_name, last_name, type, age\\) VALUES \\(\\$1, \\$2, \\$3, \\$4\\) RETURNING id").
		WithArgs("John", "Doe", "student", 25).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	person := CompletePerson{
		FirstName: "John",
		LastName:  "Doe",
		Type:      "student",
		Age:       25,
	}
	body, err := json.Marshal(person)
	assert.NoError(t, err)

	req, err := http.NewRequest("POST", "/person", bytes.NewBuffer(body))
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler.CreatePerson(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)
	var response map[string]uint
	err = json.NewDecoder(rr.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, uint(1), response["id"])
}

func TestDeletePerson(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	handler := &RequestHandler{DB: db}

	// Mock the database response
	mock.ExpectQuery("SELECT id FROM person WHERE first_name \\|\\| ' ' \\|\\| last_name = \\$1").
		WithArgs("John Doe").WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectExec("DELETE FROM person_course WHERE person_id = \\$1").
		WithArgs(1).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("DELETE FROM person WHERE id = \\$1").
		WithArgs(1).WillReturnResult(sqlmock.NewResult(1, 1))

	req, err := http.NewRequest("DELETE", "/person/John Doe", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("name", "John Doe")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	handler.DeletePerson(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	var response map[string]string
	err = json.NewDecoder(rr.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, "Person deleted successfully", response["message"])
}
