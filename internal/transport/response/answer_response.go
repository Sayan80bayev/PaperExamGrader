package response

type AnswerResponse struct {
	ID     uint    `json:"id"`
	ExamID uint    `json:"exam_id"`
	PdfURL string  `json:"pdf_url"`
	Grade  float32 `json:"grade,omitempty"`
	Images []Image ` json:"images"`
}
