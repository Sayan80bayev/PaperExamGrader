// FILEPATH: /Users/sayanseksenbaev/Programming/PaperExamGrader/internal/service/class_service.go

package service

import (
	"PaperExamGrader/internal/model"
	"PaperExamGrader/internal/repository"
	"gorm.io/gorm"
)

type ClassService struct {
	repo *repository.ClassRepository
}

func NewClassService(db *gorm.DB) *ClassService {
	return &ClassService{
		repo: repository.NewClassRepository(db),
	}
}

func (s *ClassService) Create(class *model.Class) error {
	return s.repo.Create(class)
}

func (s *ClassService) GetByID(id uint) (*model.Class, error) {
	return s.repo.GetByID(id)
}

func (s *ClassService) Update(class *model.Class) error {
	return s.repo.Update(class)
}

func (s *ClassService) Delete(id uint) error {
	return s.repo.Delete(id)
}

func (s *ClassService) List() ([]model.Class, error) {
	return s.repo.List()
}
