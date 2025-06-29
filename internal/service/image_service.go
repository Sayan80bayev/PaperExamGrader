// FILEPATH: /Users/sayanseksenbaev/Programming/PaperExamGrader/internal/service/image_service.go

package service

import (
	"PaperExamGrader/internal/model"
	"PaperExamGrader/internal/repository"
	"gorm.io/gorm"
)

type ImageService struct {
	repo *repository.ImageRepository
}

func NewImageService(db *gorm.DB) *ImageService {
	return &ImageService{
		repo: repository.NewImageRepository(db),
	}
}

func (s *ImageService) Create(image *model.Image) error {
	return s.repo.Create(image)
}

func (s *ImageService) GetByID(id uint) (*model.Image, error) {
	return s.repo.GetByID(id)
}

func (s *ImageService) Update(image *model.Image) error {
	return s.repo.Update(image)
}

func (s *ImageService) Delete(id uint) error {
	return s.repo.Delete(id)
}

func (s *ImageService) List() ([]model.Image, error) {
	return s.repo.List()
}
