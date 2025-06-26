package model

import "gorm.io/gorm"

type Class struct {
	gorm.Model
	InstructorID uint
	Instructor   User
	CRN          string
}
