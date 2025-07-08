package model

import "gorm.io/gorm"

type Answer struct {
	gorm.Model
	StudentID uint    `json:"student_id"`
	ExamID    uint    `json:"exam_id"`
	Exam      Exam    `json:"exam"`
	Grade     float32 `json:"grade,omitempty"`
	PdfURL    string  `json:"pdf_url"`
}
