package repository

import (
	"PaperExamGrader/internal/model"
	"gorm.io/gorm"
	"sync"
)

type AnswerRepository struct {
	db *gorm.DB
}

var answerRepository *AnswerRepository
var answerRepoOnce sync.Once

func GetAnswerRepository(db *gorm.DB) *AnswerRepository {
	answerRepoOnce.Do(func() {
		answerRepository = &AnswerRepository{db: db}
	})
	return answerRepository
}

// ✅ Создание нового ответа (PDF уже загружен, передаётся URL)
func (r *AnswerRepository) CreateAnswer(examID uint, pdfURL string) (*model.Answer, error) {
	answer := &model.Answer{
		ExamID: examID,
		PdfURL: pdfURL,
	}
	if err := r.db.Create(answer).Error; err != nil {
		return nil, err
	}
	return answer, nil
}

// ✅ Получение всех ответов по экзамену
func (r *AnswerRepository) GetByExamID(examID uint) ([]model.Answer, error) {
	var answers []model.Answer
	if err := r.db.Where("exam_id = ?", examID).Find(&answers).Error; err != nil {
		return nil, err
	}
	return answers, nil
}

func (r *AnswerRepository) GetWithImagesByExamID(examID uint) ([]model.Answer, error) {
	var answers []model.Answer
	if err := r.db.Preload("Images").Where("exam_id = ?", examID).Find(&answers).Error; err != nil {
		return nil, err
	}
	return answers, nil
}

// ✅ Получение одного ответа по ID
func (r *AnswerRepository) GetByID(id uint) (*model.Answer, error) {
	var answer model.Answer
	if err := r.db.First(&answer, id).Error; err != nil {
		return nil, err
	}
	return &answer, nil
}

func (r *AnswerRepository) GetWithImagesByID(id uint) (*model.Answer, error) {
	var answer model.Answer
	if err := r.db.Preload("Images").First(&answer, id).Error; err != nil {
		return nil, err
	}
	return &answer, nil
}

// ✅ Обновление оценки
func (r *AnswerRepository) UpdateGrade(id uint, grade float32) error {
	return r.db.Model(&model.Answer{}).Where("id = ?", id).Update("grade", grade).Error
}

// ✅ Удаление записи из базы (без удаления файла)
func (r *AnswerRepository) Delete(id uint) error {
	return r.db.Delete(&model.Answer{}, id).Error
}
