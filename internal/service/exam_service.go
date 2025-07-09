package service

import (
	"PaperExamGrader/internal/model"
	"PaperExamGrader/internal/repository"
	"PaperExamGrader/internal/transport/request"
	"PaperExamGrader/pkg/logging"
	"errors"
	"gorm.io/gorm"
)

type ExamService struct {
	repo *repository.ExamRepository
}

func NewExamService(db *gorm.DB) *ExamService {
	return &ExamService{
		repo: repository.NewExamRepository(db),
	}
}

func (s *ExamService) Create(req request.ExamRequest, instructorId uint) error {
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

func (s *ExamService) GetByID(id uint) (*model.Exam, error) {
	return s.repo.GetByID(id)
}

func (s *ExamService) Update(req request.ExamRequest, examId uint, instructorId uint) error {
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

func (s *ExamService) Delete(id uint) error {
	return s.repo.Delete(id)
}

func (s *ExamService) List() ([]model.Exam, error) {
	return s.repo.List()
}
