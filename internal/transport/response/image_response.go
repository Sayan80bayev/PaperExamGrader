package response

type Image struct {
	ID       uint   `json:"id"`
	AnswerID uint   `json:"answer_id"`
	URL      string `json:"url"`
}
