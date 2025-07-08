package model

import "gorm.io/gorm"

type Exam struct {
	gorm.Model   `json:"gorm_._model"`
	InstructorID uint   `json:"instructor_id"`
	CRN          string `json:"crn"`
	Date         string `json:"date"`
}
