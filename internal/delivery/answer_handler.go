package delivery

import (
	"PaperExamGrader/internal/service"
	"PaperExamGrader/internal/transport/response"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

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
