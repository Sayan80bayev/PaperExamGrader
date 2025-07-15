package response

import "time"

type BBoxTemplateResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	ExamID    uint      `json:"exam_id"`
	CreatedAt time.Time `json:"created_at"`
}
