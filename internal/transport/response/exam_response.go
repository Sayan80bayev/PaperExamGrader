package response

import "gorm.io/gorm"

type Exam struct {
	gorm.Model `json:"gorm_._model"`
	CRN        string `json:"crn"`
	Date       string `json:"date"`
}
