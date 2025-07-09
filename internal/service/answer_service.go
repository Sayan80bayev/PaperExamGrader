package service

import (
	"PaperExamGrader/internal/model"
	"PaperExamGrader/internal/repository"
	"PaperExamGrader/internal/storage"
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
	// Сначала загрузим файл в MinIO
	url, err := s.fileStorage.UploadFile(file, header)
	if err != nil {
		return nil, err
	}

	// Затем создадим ответ в базе с ссылкой
	return s.answerRepo.CreateAnswer(examID, url)
}

// ✅ Получить все ответы по экзамену
func (s *AnswerService) GetAnswersByExam(examID uint) ([]model.Answer, error) {
	return s.answerRepo.GetByExamID(examID)
}

// ✅ Получить конкретный ответ по ID
func (s *AnswerService) GetAnswerByID(id uint) (*model.Answer, error) {
	return s.answerRepo.GetByID(id)
}

// ✅ Обновить оценку ответа
func (s *AnswerService) UpdateAnswerGrade(id uint, grade float32) error {
	return s.answerRepo.UpdateGrade(id, grade)
}

// ✅ Удалить ответ (включая удаление файла)
func (s *AnswerService) DeleteAnswer(id uint) error {
	// Сначала получим ответ
	answer, err := s.answerRepo.GetByID(id)
	if err != nil {
		return err
	}

	// Удалим PDF из MinIO
	if err := s.fileStorage.DeleteFileByURL(answer.PdfURL); err != nil {
		return err
	}

	// Удалим запись из базы
	return s.answerRepo.Delete(id)
}
