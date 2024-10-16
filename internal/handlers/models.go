package handlers

import "database/sql"

type Course struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type Person struct {
	ID        uint   `json: "id"`
	FirstName string `json: "first_name"`
	LastName  string `json: "last_name"`
	Type      string `json: "type"` //only 'student' or 'professor'
	Age       uint   `json: "age"`
}

type CompletePerson struct {
	ID        uint   `json: "id"`
	FirstName string `json: "first_name"`
	LastName  string `json: "last_name"`
	Type      string `json: "type"` //only 'student' or 'professor'
	Age       uint   `json: "age"`
	Courses   []uint `json: "courses"`
}

type PersonCourse struct {
	PersonID uint `json:"person_id"`
	CourseID uint `json:"person_id"`
}

// set table name
func (Course) TableName() string {
	return "course"
}

func (Person) TableName() string {
	return "person"
}

func (PersonCourse) TableName() string {
	return "person_course"
}

type RequestHandler struct {
	DB *sql.DB
}
