package request

type ExamRequest struct {
	CRN  string `json:"crn" binding:"required"`
	Date string `json:"date" binding:"required"`
}
