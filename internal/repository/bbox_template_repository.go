package repository

import (
	"PaperExamGrader/internal/model"

	"gorm.io/gorm"
)

type BBoxTemplateRepository interface {
	CreateTemplate(template *model.BBoxTemplate) error
	GetTemplatesByExamID(examID uint) ([]model.BBoxTemplate, error)
	GetByID(id uint) (*model.BBoxTemplate, error)
	DeleteTemplate(id uint) error
}

type bboxTemplateRepository struct {
	db *gorm.DB
}

func GetBBoxTemplateRepository(db *gorm.DB) BBoxTemplateRepository {
	return &bboxTemplateRepository{db}
}

func (r *bboxTemplateRepository) CreateTemplate(template *model.BBoxTemplate) error {
	return r.db.Create(template).Error
}

func (r *bboxTemplateRepository) GetTemplatesByExamID(examID uint) ([]model.BBoxTemplate, error) {
	var templates []model.BBoxTemplate
	if err := r.db.Preload("BBoxes").Where("exam_id = ?", examID).Find(&templates).Error; err != nil {
		return nil, err
	}
	return templates, nil
}

func (r *bboxTemplateRepository) GetByID(id uint) (*model.BBoxTemplate, error) {
	var template model.BBoxTemplate
	if err := r.db.First(&template, id).Error; err != nil {
		return nil, err
	}
	return &template, nil
}

func (r *bboxTemplateRepository) DeleteTemplate(id uint) error {
	return r.db.Delete(&model.BBoxTemplate{}, id).Error
}
