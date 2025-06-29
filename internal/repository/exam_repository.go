// FILEPATH: /Users/sayanseksenbaev/Programming/PaperExamGrader/internal/repository/exam_repository.go

package repository

import (
	"sync"

	"PaperExamGrader/internal/model"
	"gorm.io/gorm"
)

type ExamRepository struct {
	db *gorm.DB
}

var examRepoInstance *ExamRepository
var examRepoOnce sync.Once

func NewExamRepository(db *gorm.DB) *ExamRepository {
	examRepoOnce.Do(func() {
		examRepoInstance = &ExamRepository{db: db}
	})
	return examRepoInstance
}

func (r *ExamRepository) Create(exam *model.Exam) error {
	return r.db.Create(exam).Error
}

func (r *ExamRepository) GetByID(id uint) (*model.Exam, error) {
	var exam model.Exam
	err := r.db.First(&exam, id).Error
	return &exam, err
}

func (r *ExamRepository) Update(exam *model.Exam) error {
	return r.db.Save(exam).Error
}

func (r *ExamRepository) Delete(id uint) error {
	return r.db.Delete(&model.Exam{}, id).Error
}

func (r *ExamRepository) List() ([]model.Exam, error) {
	var exams []model.Exam
	err := r.db.Find(&exams).Error
	return exams, err
}
