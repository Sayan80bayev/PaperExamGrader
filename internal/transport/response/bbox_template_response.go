package response

import (
	"PaperExamGrader/internal/model"
	"time"
)

type BBoxTemplate struct {
	ID        uint               `json:"id"`
	Name      string             `json:"name"`
	ExamID    uint               `json:"exam_id"`
	CreatedAt time.Time          `json:"created_at"`
	BBoxes    []model.BBoxMetaDB `json:"bboxes"`
}
