package service

import (
	"PaperExamGrader/internal/model"
	"PaperExamGrader/internal/repository"
)

type BBoxService interface {
	Create(meta *model.BBoxMetaDB) error
	GetByID(id uint) (*model.BBoxMetaDB, error)
	GetAllByExamID(id uint) ([]model.BBoxMetaDB, error)
	Update(meta *model.BBoxMetaDB) error
	Delete(id uint) error
}

type bboxService struct {
	repo     repository.BBoxRepository
	examRepo *repository.ExamRepository
}

func NewBBoxService(repo repository.BBoxRepository, examRepository *repository.ExamRepository) BBoxService {
	return &bboxService{repo: repo, examRepo: examRepository}
}

func (s *bboxService) Create(meta *model.BBoxMetaDB) error {
	return s.repo.Create(meta)
}

func (s *bboxService) GetByID(id uint) (*model.BBoxMetaDB, error) {
	return s.repo.GetByID(id)
}

func (s *bboxService) GetAllByExamID(id uint) ([]model.BBoxMetaDB, error) {
	_, err := s.examRepo.GetByID(id)
	if err != nil {
		return []model.BBoxMetaDB{}, err
	}

	return s.repo.GetAllByExamID(id)
}

func (s *bboxService) Update(meta *model.BBoxMetaDB) error {
	return s.repo.Update(meta)
}

func (s *bboxService) Delete(id uint) error {
	return s.repo.Delete(id)
}
