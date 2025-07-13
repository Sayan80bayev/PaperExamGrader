package delivery

import (
	"PaperExamGrader/internal/service"
	"PaperExamGrader/internal/transport/response"
	"archive/zip"
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"strconv"
	"strings"
)

// 👇 Custom type to implement multipart.File
type readSeekCloser struct {
	*bytes.Reader
}

func (r *readSeekCloser) Close() error {
	return nil
}

type AnswerHandler struct {
	answerService *service.AnswerService
}

func NewAnswerHandler(s *service.AnswerService) *AnswerHandler {
	return &AnswerHandler{answerService: s}
}

// ✅ POST /answers/upload?exam_id=1
func (h *AnswerHandler) Upload(c *gin.Context) {
	examIDStr := c.Query("exam_id")
	examID, err := strconv.Atoi(examIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid exam_id"})
		return
	}

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file upload error"})
		return
	}
	defer file.Close()

	answer, err := h.answerService.UploadAnswer(file, header, uint(examID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to upload answer"})
		return
	}

	resp := response.AnswerResponse{
		ID:     answer.ID,
		ExamID: answer.ExamID,
		PdfURL: answer.PdfURL,
		Grade:  answer.Grade,
	}
	c.JSON(http.StatusOK, resp)
}

// ✅ POST /answers/upload_zip?exam_id=1
func (h *AnswerHandler) UploadFromZip(c *gin.Context) {
	examIDStr := c.Query("exam_id")
	examID, err := strconv.Atoi(examIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid exam_id"})
		return
	}

	zipFile, _, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to get zip file"})
		return
	}
	defer zipFile.Close()

	// Read ZIP into memory
	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, zipFile); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read zip file"})
		return
	}

	zipReader, err := zip.NewReader(bytes.NewReader(buf.Bytes()), int64(buf.Len()))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid zip file"})
		return
	}

	var responses []response.AnswerResponse

	for _, file := range zipReader.File {
		if file.FileInfo().IsDir() {
			continue // skip directories
		}

		if strings.HasPrefix(file.Name, "__MACOSX/") || strings.HasPrefix(file.Name, "._") || strings.Contains(file.Name, "/._") {
			continue
		}

		f, err := file.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to open file in zip: %s", file.Name)})
			return
		}

		// Copy file content to memory
		fileCopy := bytes.NewBuffer(nil)
		if _, err := io.Copy(fileCopy, f); err != nil {
			_ = f.Close()
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to read file: %s", file.Name)})
			return
		}
		_ = f.Close()

		// Wrap bytes.Reader to satisfy multipart.File
		reader := bytes.NewReader(fileCopy.Bytes())
		multipartFile := &readSeekCloser{Reader: reader}

		// Create fake multipart.FileHeader
		fileHeader := &multipart.FileHeader{
			Filename: file.Name,
			Size:     int64(file.UncompressedSize64),
			Header:   make(textproto.MIMEHeader),
		}
		fileHeader.Header.Set("Content-Disposition", fmt.Sprintf(`form-data; name="file"; filename="%s"`, file.Name))
		fileHeader.Header.Set("Content-Type", "application/octet-stream")

		// Upload each file
		answer, err := h.answerService.UploadAnswer(multipartFile, fileHeader, uint(examID))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to upload: %s", file.Name)})
			return
		}

		resp := response.AnswerResponse{
			ID:     answer.ID,
			ExamID: answer.ExamID,
			PdfURL: answer.PdfURL,
			Grade:  answer.Grade,
		}
		responses = append(responses, resp)
	}

	c.JSON(http.StatusOK, responses)
}

// ✅ GET /answers/:id
func (h *AnswerHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	answer, err := h.answerService.GetAnswerByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	resp := response.AnswerResponse{
		ID:     answer.ID,
		ExamID: answer.ExamID,
		PdfURL: answer.PdfURL,
		Grade:  answer.Grade,
	}
	c.JSON(http.StatusOK, resp)
}

// ✅ GET /answers/exam/:exam_id
func (h *AnswerHandler) GetByExamID(c *gin.Context) {
	examIDStr := c.Param("exam_id")
	examID, err := strconv.Atoi(examIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid exam_id"})
		return
	}

	answers, err := h.answerService.GetAnswersByExam(uint(examID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve answers"})
		return
	}

	var resp []response.AnswerResponse
	for _, ans := range answers {
		resp = append(resp, response.AnswerResponse{
			ID:     ans.ID,
			ExamID: ans.ExamID,
			PdfURL: ans.PdfURL,
			Grade:  ans.Grade,
		})
	}
	c.JSON(http.StatusOK, resp)
}

// ✅ PUT /answers/:id/grade
func (h *AnswerHandler) UpdateGrade(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req struct {
		Grade float32 `json:"grade"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	if err := h.answerService.UpdateAnswerGrade(uint(id), req.Grade); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update grade"})
		return
	}

	c.Status(http.StatusOK)
}

// ✅ DELETE /answers/:id
func (h *AnswerHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.answerService.DeleteAnswer(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "delete failed"})
		return
	}

	c.Status(http.StatusNoContent)
}
