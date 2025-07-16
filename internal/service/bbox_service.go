package service

import (
	"PaperExamGrader/internal/model"
	"PaperExamGrader/internal/repository"
	"PaperExamGrader/internal/transport/request"
	"PaperExamGrader/internal/transport/response"
	"errors"

	"gorm.io/gorm"
)

type BBoxService struct {
	db           *gorm.DB
	bboxRepo     repository.BBoxRepository
	templateRepo repository.BBoxTemplateRepository
	examRepo     *repository.ExamRepository
}

func NewBBoxService(db *gorm.DB, bboxRepo repository.BBoxRepository, templateRepo repository.BBoxTemplateRepository, examRepository *repository.ExamRepository) *BBoxService {
	return &BBoxService{
		db:           db,
		bboxRepo:     bboxRepo,
		templateRepo: templateRepo,
		examRepo:     examRepository,
	}
}

func (s *BBoxService) CreateTemplateWithBBoxes(req request.CreateBBoxTemplateRequest, instructorId uint) (*response.BBoxTemplate, error) {
	exam, err := s.examRepo.GetByID(req.ExamID)
	if err != nil {
		return nil, err
	}
	if exam.InstructorID != instructorId {
		return nil, errors.New("you are not the owner of this exam")
	}

	template := &model.BBoxTemplate{
		ExamID: req.ExamID,
		Name:   req.Name,
	}

	err = s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(template).Error; err != nil {
			return err
		}

		for i := range req.BBoxes {
			req.BBoxes[i].TemplateID = template.ID
		}

		if err := tx.Create(&req.BBoxes).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &response.BBoxTemplate{
		ID:        template.ID,
		Name:      template.Name,
		ExamID:    template.ExamID,
		CreatedAt: template.CreatedAt,
	}, nil
}

func (s *BBoxService) GetTemplatesByExamID(examID, instructorID uint) ([]response.BBoxTemplate, error) {
	exam, err := s.examRepo.GetByID(examID)
	if err != nil {
		return nil, err
	}
	if exam.InstructorID != instructorID {
		return nil, errors.New("you are not the owner of this exam")
	}

	templates, err := s.templateRepo.GetTemplatesByExamID(examID)
	if err != nil {
		return nil, err
	}

	res := make([]response.BBoxTemplate, len(templates))
	for i, t := range templates {
		res[i] = response.BBoxTemplate{
			ID:        t.ID,
			Name:      t.Name,
			ExamID:    t.ExamID,
			CreatedAt: t.CreatedAt,
			BBoxes:    t.BBoxes,
		}
	}
	return res, nil
}

func (s *BBoxService) DeleteTemplate(id, instructorID uint) error {
	template, err := s.templateRepo.GetByID(id)
	if err != nil {
		return err
	}
	exam, err := s.examRepo.GetByID(template.ExamID)
	if err != nil {
		return err
	}
	if exam.InstructorID != instructorID {
		return errors.New("you are not the owner of this exam")
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("template_id = ?", id).Delete(&model.BBoxMetaDB{}).Error; err != nil {
			return err
		}
		if err := tx.Delete(&model.BBoxTemplate{}, id).Error; err != nil {
			return err
		}
		return nil
	})
}
