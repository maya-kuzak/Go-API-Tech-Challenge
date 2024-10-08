// create person struct
package models

type Person struct {
	ID        uint   `gorm:"primaryKey;autoIncrement" json: "id"`
	FirstName string `gorm:"not null" json: "first_name"`
	LastName  string `gorm:"not null" json: "last_name"`
	Type      string `gorm:"not null" validate:"oneof=professor student" json: "type"` //only 'student' or 'professor'
	Age       uint   `gorm:"not null" json: "age"`
}

// set table name
func (Person) TableName() string {
	return "person"
}

// type PersonCourse struct {
// 	Person
// 	Courses []int `gorm:"many2many:person_course" json:"courses"`
// }
