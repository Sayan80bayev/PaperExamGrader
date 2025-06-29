// FILEPATH: /Users/sayanseksenbaev/Programming/PaperExamGrader/internal/repository/class_repository.go

package repository

import (
	"sync"

	"PaperExamGrader/internal/model"
	"gorm.io/gorm"
)

type ClassRepository struct {
	db *gorm.DB
}

var classRepoInstance *ClassRepository
var classRepoOnce sync.Once

func NewClassRepository(db *gorm.DB) *ClassRepository {
	classRepoOnce.Do(func() {
		classRepoInstance = &ClassRepository{db: db}
	})
	return classRepoInstance
}

func (r *ClassRepository) Create(class *model.Class) error {
	return r.db.Create(class).Error
}

func (r *ClassRepository) GetByID(id uint) (*model.Class, error) {
	var class model.Class
	err := r.db.First(&class, id).Error
	return &class, err
}

func (r *ClassRepository) Update(class *model.Class) error {
	return r.db.Save(class).Error
}

func (r *ClassRepository) Delete(id uint) error {
	return r.db.Delete(&model.Class{}, id).Error
}

func (r *ClassRepository) List() ([]model.Class, error) {
	var classes []model.Class
	err := r.db.Find(&classes).Error
	return classes, err
}
