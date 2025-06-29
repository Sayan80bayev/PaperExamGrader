// FILEPATH: /Users/sayanseksenbaev/Programming/PaperExamGrader/internal/service/user_service.go

package service

import (
	"PaperExamGrader/internal/model"
	"PaperExamGrader/internal/repository"
	"gorm.io/gorm"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{
		repo: repository.NewUserRepository(db),
	}
}

func (s *UserService) Create(user *model.User) error {
	return s.repo.Create(user)
}

func (s *UserService) GetByID(id uint) (*model.User, error) {
	return s.repo.GetByID(id)
}

func (s *UserService) Update(user *model.User) error {
	return s.repo.Update(user)
}

func (s *UserService) Delete(id uint) error {
	return s.repo.Delete(id)
}

func (s *UserService) List() ([]model.User, error) {
	return s.repo.List()
}
