package request

type AnswerRequest struct {
	ExamID uint `form:"exam_id" binding:"required"`
	// PDF file обрабатывается через multipart/form-data
}
