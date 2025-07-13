package service

import (
	"PaperExamGrader/internal/model"
	"PaperExamGrader/internal/repository"
	"PaperExamGrader/internal/transport/request"
	"PaperExamGrader/internal/transport/response"
	"PaperExamGrader/pkg/logging"
	"errors"
	"gorm.io/gorm"
)

type ExamService struct {
	repo *repository.ExamRepository
}

func NewExamService(db *gorm.DB) *ExamService {
	return &ExamService{
		repo: repository.GetExamRepository(db),
	}
}

func (s *ExamService) Create(req request.Exam, instructorId uint) error {
	if len(req.CRN) < 2 || len(req.CRN) > 6 {
		return errors.New("CRN must be at least 2 and at most 6 characters long")
	}

	exam := &model.Exam{
		CRN:          req.CRN,
		Date:         req.Date,
		InstructorID: instructorId,
	}

	return s.repo.Create(exam)
}

func (s *ExamService) GetByID(id uint) (*response.Exam, error) {
	examModel, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	examResponse := &response.Exam{
		Model: examModel.Model,
		CRN:   examModel.CRN,
		Date:  examModel.Date,
	}

	return examResponse, nil
}

func (s *ExamService) Update(req request.Exam, examId uint, instructorId uint) error {
	var log = logging.GetLogger()

	exam, err := s.repo.GetByID(examId)
	if err != nil {
		log.Warnf("Failed to get exam by id: %d", examId)
		return err
	}

	if exam.InstructorID != instructorId {
		return errors.New("You are not the owner of this exam ")
	}

	exam.CRN = req.CRN
	exam.Date = req.Date

	return s.repo.Update(exam)
}

func (s *ExamService) Delete(id uint, instructorId uint) error {
	log := logging.GetLogger()

	exam, err := s.repo.GetByID(id)
	if exam == nil || err != nil {
		log.Warnf("Failed to get exam by id: %d", id)
		return errors.New("Could not found the exam ")
	} else if exam.InstructorID != instructorId {
		return errors.New("You are not the owner of this exam ")
	}

	return s.repo.Delete(id)
}

func (s *ExamService) List(id uint) ([]response.Exam, error) {
	examModels, err := s.repo.GetAllByUserID(id)
	if err != nil {
		return nil, err
	}

	examResponses := make([]response.Exam, len(examModels))
	for i, exam := range examModels {
		examResponses[i] = response.Exam{
			Model: exam.Model,
			CRN:   exam.CRN,
			Date:  exam.Date,
		}
	}

	return examResponses, nil
}
