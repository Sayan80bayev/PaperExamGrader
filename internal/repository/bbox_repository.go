package repository

import (
	"PaperExamGrader/internal/model"

	"gorm.io/gorm"
)

type BBoxRepository interface {
	CreateMany(bboxes []model.BBoxMetaDB) error
}

type bboxRepository struct {
	db *gorm.DB
}

func GetBBoxRepository(db *gorm.DB) BBoxRepository {
	return &bboxRepository{db}
}

func (r *bboxRepository) CreateMany(bboxes []model.BBoxMetaDB) error {
	return r.db.Create(&bboxes).Error
}
