package model

import "gorm.io/gorm"

type Answer struct {
	gorm.Model
	StudentID uint    `json:"student_id"`
	ClassID   uint    `json:"class_id"`
	Class     Exam    `json:"class"`
	Grade     float32 `json:"grade,omitempty"`
	PdfURL    string  `json:"pdf_url"`
}
