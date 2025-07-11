package repository

import (
	"PaperExamGrader/internal/model"
	"sync"

	"gorm.io/gorm"
)

type BBoxRepository interface {
	Create(bbox *model.BBoxMetaDB) error
	GetByID(id uint) (*model.BBoxMetaDB, error)
	GetAll() ([]model.BBoxMetaDB, error)
	Update(bbox *model.BBoxMetaDB) error
	Delete(id uint) error
}

type bboxRepo struct {
	db *gorm.DB
}

var (
	repoInstance BBoxRepository
	once         sync.Once
)

// GetBBoxRepository returns a singleton instance of the repository
func GetBBoxRepository(db *gorm.DB) BBoxRepository {
	once.Do(func() {
		repoInstance = &bboxRepo{db: db}
	})
	return repoInstance
}

func (r *bboxRepo) Create(bbox *model.BBoxMetaDB) error {
	return r.db.Create(bbox).Error
}

func (r *bboxRepo) GetByID(id uint) (*model.BBoxMetaDB, error) {
	var bbox model.BBoxMetaDB
	err := r.db.First(&bbox, id).Error
	return &bbox, err
}

func (r *bboxRepo) GetAll() ([]model.BBoxMetaDB, error) {
	var bboxes []model.BBoxMetaDB
	err := r.db.Find(&bboxes).Error
	return bboxes, err
}

func (r *bboxRepo) Update(bbox *model.BBoxMetaDB) error {
	return r.db.Save(bbox).Error
}

func (r *bboxRepo) Delete(id uint) error {
	return r.db.Delete(&model.BBoxMetaDB{}, id).Error
}
