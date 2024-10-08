// create course struct
package models

type Course struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

// set table name
func (Course) TableName() string {
	return "course"
}
