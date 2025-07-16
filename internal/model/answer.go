package model

import (
	"gorm.io/gorm"
)

type Answer struct {
	gorm.Model
	ExamID uint    `json:"exam_id"`
	Exam   Exam    `json:"exam"`
	Grade  float32 `json:"grade,omitempty"`
	PdfURL string  `json:"pdf_url"`
	Images []Image `gorm:"foreignKey:AnswerID" json:"images"`
}
