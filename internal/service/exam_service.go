// FILEPATH: /Users/sayanseksenbaev/Programming/PaperExamGrader/internal/service/exam_service.go

package service

import (
	"PaperExamGrader/internal/model"
	"PaperExamGrader/internal/repository"
	"gorm.io/gorm"
)

type ExamService struct {
	repo *repository.ExamRepository
}

func NewExamService(db *gorm.DB) *ExamService {
	return &ExamService{
		repo: repository.NewExamRepository(db),
	}
}

func (s *ExamService) Create(exam *model.Exam) error {
	return s.repo.Create(exam)
}

func (s *ExamService) GetByID(id uint) (*model.Exam, error) {
	return s.repo.GetByID(id)
}

func (s *ExamService) Update(exam *model.Exam) error {
	return s.repo.Update(exam)
}

func (s *ExamService) Delete(id uint) error {
	return s.repo.Delete(id)
}

func (s *ExamService) List() ([]model.Exam, error) {
	return s.repo.List()
}
