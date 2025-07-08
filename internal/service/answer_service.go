// FILEPATH: /Users/sayanseksenbaev/Programming/PaperAnswerGrader/internal/service/answer_service.go

package service

import (
	"PaperExamGrader/internal/model"
	"PaperExamGrader/internal/repository"
	"gorm.io/gorm"
)

type AnswerService struct {
	repo *repository.AnswerRepository
}

func NewAnswerService(db *gorm.DB) *AnswerService {
	return &AnswerService{
		repo: repository.NewAnswerRepository(db),
	}
}

func (s *AnswerService) Create(answer *model.Answer) error {
	return s.repo.Create(answer)
}

func (s *AnswerService) GetByID(id uint) (*model.Answer, error) {
	return s.repo.GetByID(id)
}

func (s *AnswerService) Update(answer *model.Answer) error {
	return s.repo.Update(answer)
}

func (s *AnswerService) Delete(id uint) error {
	return s.repo.Delete(id)
}

func (s *AnswerService) List() ([]model.Answer, error) {
	return s.repo.List()
}
