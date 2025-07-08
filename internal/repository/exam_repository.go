// FILEPATH: /Users/sayanseksenbaev/Programming/PaperExamGrader/internal/repository/class_repository.go

package repository

import (
	"sync"

	"PaperExamGrader/internal/model"
	"gorm.io/gorm"
)

type ExamRepository struct {
	db *gorm.DB
}

var examRepository *ExamRepository
var examRepoOnce sync.Once

func NewExamRepository(db *gorm.DB) *ExamRepository {
	examRepoOnce.Do(func() {
		examRepository = &ExamRepository{db: db}
	})
	return examRepository
}

func (r *ExamRepository) Create(class *model.Exam) error {
	return r.db.Create(class).Error
}

func (r *ExamRepository) GetByID(id uint) (*model.Exam, error) {
	var class model.Exam
	err := r.db.First(&class, id).Error
	return &class, err
}

func (r *ExamRepository) Update(class *model.Exam) error {
	return r.db.Save(class).Error
}

func (r *ExamRepository) Delete(id uint) error {
	return r.db.Delete(&model.Exam{}, id).Error
}

func (r *ExamRepository) List() ([]model.Exam, error) {
	var classes []model.Exam
	err := r.db.Find(&classes).Error
	return classes, err
}
