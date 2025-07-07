package model

import "gorm.io/gorm"

type Class struct {
	gorm.Model   `json:"gorm_._model"`
	InstructorID uint   `json:"instructor_id"`
	CRN          string `json:"crn"`
}
