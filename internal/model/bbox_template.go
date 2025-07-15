package model

import "time"

type BBoxTemplate struct {
	ID        uint         `gorm:"primaryKey" json:"id"`
	ExamID    uint         `json:"exam_id"`
	Name      string       `json:"name"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
	BBoxes    []BBoxMetaDB `gorm:"foreignKey:TemplateID"`
}
