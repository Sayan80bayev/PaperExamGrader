package service

import (
	"PaperExamGrader/internal/model"
	"PaperExamGrader/internal/repository"
	"PaperExamGrader/internal/storage"
	"PaperExamGrader/internal/transport/response"
	"mime/multipart"
)

type AnswerService struct {
	answerRepo  *repository.AnswerRepository
	fileStorage storage.FileStorage
}

func NewAnswerService(repo *repository.AnswerRepository, fs storage.FileStorage) *AnswerService {
	return &AnswerService{
		answerRepo:  repo,
		fileStorage: fs,
	}
}

// ✅ Загрузка нового ответа
func (s *AnswerService) UploadAnswer(file multipart.File, header *multipart.FileHeader, examID uint) (*model.Answer, error) {
	url, err := s.fileStorage.UploadFile(file, header)
	if err != nil {
		return nil, err
	}
	return s.answerRepo.CreateAnswer(examID, url)
}

// ✅ Получить все ответы по экзамену (с отображением)
func (s *AnswerService) GetAnswersByExam(examID uint) ([]response.AnswerResponse, error) {
	answers, err := s.answerRepo.GetByExamID(examID)
	if err != nil {
		return nil, err
	}
	return mapAnswersToResponses(answers), nil
}

// ✅ Получить ответы с изображениями по экзамену (с отображением)
func (s *AnswerService) GetAnswersWithImagesByExam(examID uint) ([]response.AnswerResponse, error) {
	answers, err := s.answerRepo.GetWithImagesByExamID(examID)
	if err != nil {
		return nil, err
	}
	return mapAnswersToResponses(answers), nil
}

// ✅ Получить конкретный ответ по ID (с отображением)
func (s *AnswerService) GetAnswerByID(id uint) (*response.AnswerResponse, error) {
	answer, err := s.answerRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	resp := mapAnswerToResponse(*answer)
	return &resp, nil
}

func (s *AnswerService) GetAnswerWithImagesByID(id uint) (*response.AnswerResponse, error) {
	answer, err := s.answerRepo.GetWithImagesByID(id)
	if err != nil {
		return nil, err
	}
	resp := mapAnswerToResponse(*answer)
	return &resp, nil
}

// ✅ Обновить оценку ответа
func (s *AnswerService) UpdateAnswerGrade(id uint, grade float32) error {
	return s.answerRepo.UpdateGrade(id, grade)
}

// ✅ Удалить ответ (включая удаление файла)
func (s *AnswerService) DeleteAnswer(id uint) error {
	answer, err := s.answerRepo.GetByID(id)
	if err != nil {
		return err
	}

	if err := s.fileStorage.DeleteFileByURL(answer.PdfURL); err != nil {
		return err
	}

	return s.answerRepo.Delete(id)
}

// ✅ Map single Answer model to response
func mapAnswerToResponse(a model.Answer) response.AnswerResponse {
	images := make([]response.Image, len(a.Images))
	for i, img := range a.Images {
		images[i] = response.Image{
			ID:       img.ID,
			AnswerID: img.AnswerID,
			URL:      img.URL,
		}
	}

	return response.AnswerResponse{
		ID:     a.ID,
		ExamID: a.ExamID,
		PdfURL: a.PdfURL,
		Grade:  a.Grade,
		Images: images,
	}
}

// ✅ Map list of Answer models to responses
func mapAnswersToResponses(answers []model.Answer) []response.AnswerResponse {
	responses := make([]response.AnswerResponse, len(answers))
	for i, a := range answers {
		responses[i] = mapAnswerToResponse(a)
	}
	return responses
}
