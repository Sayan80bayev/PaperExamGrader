package request

type Exam struct {
	CRN  string `json:"crn" binding:"required"`
	Date string `json:"date" binding:"required"`
}
