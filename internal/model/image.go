package model

import "gorm.io/gorm"

type Image struct {
	gorm.Model
	ExamID uint   `json:"exam_id"`
	Exam   Exam   `json:"exam"`
	URL    string `json:"url"`
}
