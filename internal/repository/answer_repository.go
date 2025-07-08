// FILEPATH: /Users/sayanseksenbaev/Programming/PaperExamGrader/internal/repository/exam_repository.go

package repository

import (
	"sync"

	"PaperExamGrader/internal/model"
	"gorm.io/gorm"
)

type AnswerRepository struct {
	db *gorm.DB
}

var answerRepository *AnswerRepository
var answerRepoOnce sync.Once

func NewAnswerRepository(db *gorm.DB) *AnswerRepository {
	answerRepoOnce.Do(func() {
		answerRepository = &AnswerRepository{db: db}
	})
	return answerRepository
}

func (r *AnswerRepository) Create(answer *model.Answer) error {
	return r.db.Create(answer).Error
}

func (r *AnswerRepository) GetByID(id uint) (*model.Answer, error) {
	var answer model.Answer
	err := r.db.First(&answer, id).Error
	return &answer, err
}

func (r *AnswerRepository) Update(answer *model.Answer) error {
	return r.db.Save(answer).Error
}

func (r *AnswerRepository) Delete(id uint) error {
	return r.db.Delete(&model.Answer{}, id).Error
}

func (r *AnswerRepository) List() ([]model.Answer, error) {
	var exams []model.Answer
	err := r.db.Find(&exams).Error
	return exams, err
}
