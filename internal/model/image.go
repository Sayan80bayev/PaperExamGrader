package model

import "gorm.io/gorm"

type Image struct {
	gorm.Model
	AnswerID uint   `json:"answer_id"`
	Answer   Answer `json:"answer"`
	URL      string `json:"url"`
}
