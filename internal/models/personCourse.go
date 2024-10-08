// many to many person course struct
package models

type PersonCourse struct {
	PersonID uint `gorm:"not null" json:"person_id"`
	CourseID uint `gorm:"not null" json:"person_id"`
}

// set table name
func (PersonCourse) TableName() string {
	return "person_course"
}
