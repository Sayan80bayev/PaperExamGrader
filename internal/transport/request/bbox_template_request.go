package request

import "PaperExamGrader/internal/model"

type CreateBBoxTemplateRequest struct {
	Name   string             `json:"name"`
	ExamID uint               `json:"exam_id"`
	BBoxes []model.BBoxMetaDB `json:"bboxes"`
}
