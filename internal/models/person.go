// create person struct
package models

type Person struct {
	ID        uint   `json: "id"`
	FirstName string `json: "first_name"`
	LastName  string `json: "last_name"`
	Type      string `json: "type"` //only 'student' or 'professor'
	Age       uint   `json: "age"`
}

// included courses
// have to be created bc separate databases, but one schema in api calls
type CompletePerson struct {
	ID        uint   `json: "id"`
	FirstName string `json: "first_name"`
	LastName  string `json: "last_name"`
	Type      string `json: "type"` //only 'student' or 'professor'
	Age       uint   `json: "age"`
	Courses   []uint `json: "courses"`
}

// set table name
func (Person) TableName() string {
	return "person"
}
