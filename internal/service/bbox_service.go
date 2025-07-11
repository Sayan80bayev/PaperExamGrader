package service

import (
	"PaperExamGrader/internal/model"
	"PaperExamGrader/internal/repository"
)

type BBoxService interface {
	Create(meta *model.BBoxMetaDB) error
	GetByID(id uint) (*model.BBoxMetaDB, error)
	GetAll() ([]model.BBoxMetaDB, error)
	Update(meta *model.BBoxMetaDB) error
	Delete(id uint) error
}

type bboxService struct {
	repo repository.BBoxRepository
}

func NewBBoxService(repo repository.BBoxRepository) BBoxService {
	return &bboxService{repo: repo}
}

func (s *bboxService) Create(meta *model.BBoxMetaDB) error {
	return s.repo.Create(meta)
}

func (s *bboxService) GetByID(id uint) (*model.BBoxMetaDB, error) {
	return s.repo.GetByID(id)
}

func (s *bboxService) GetAll() ([]model.BBoxMetaDB, error) {
	return s.repo.GetAll()
}

func (s *bboxService) Update(meta *model.BBoxMetaDB) error {
	return s.repo.Update(meta)
}

func (s *bboxService) Delete(id uint) error {
	return s.repo.Delete(id)
}
