// FILEPATH: /Users/sayanseksenbaev/Programming/PaperExamGrader/internal/repository/image_repository.go

package repository

import (
	"sync"

	"PaperExamGrader/internal/model"
	"gorm.io/gorm"
)

type ImageRepository struct {
	db *gorm.DB
}

var imageRepoInstance *ImageRepository
var imageRepoOnce sync.Once

func NewImageRepository(db *gorm.DB) *ImageRepository {
	imageRepoOnce.Do(func() {
		imageRepoInstance = &ImageRepository{db: db}
	})
	return imageRepoInstance
}

func (r *ImageRepository) Create(image *model.Image) error {
	return r.db.Create(image).Error
}

func (r *ImageRepository) GetByID(id uint) (*model.Image, error) {
	var image model.Image
	err := r.db.First(&image, id).Error
	return &image, err
}

func (r *ImageRepository) Update(image *model.Image) error {
	return r.db.Save(image).Error
}

func (r *ImageRepository) Delete(id uint) error {
	return r.db.Delete(&model.Image{}, id).Error
}

func (r *ImageRepository) List() ([]model.Image, error) {
	var images []model.Image
	err := r.db.Find(&images).Error
	return images, err
}
