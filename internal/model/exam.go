package model

import "gorm.io/gorm"

type Exam struct {
	gorm.Model
	StudentID uint    `json:"student_id"`
	ClassID   uint    `json:"class_id"`
	Class     Class   `json:"class"`
	Grade     float32 `json:"grade,omitempty"`
	PdfURL    string  `json:"pdf_url"`
}
