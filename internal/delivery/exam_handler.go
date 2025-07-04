// FILEPATH: /Users/sayanseksenbaev/Programming/PaperExamGrader/internal/handler/exam_handler.go

package delivery

import (
	"net/http"
	"strconv"

	"PaperExamGrader/internal/model"
	"PaperExamGrader/internal/service"
	"github.com/gin-gonic/gin"
)

type ExamHandler struct {
	service *service.ExamService
}

func NewExamHandler(service *service.ExamService) *ExamHandler {
	return &ExamHandler{service: service}
}

// Create creates a new exam
func (h *ExamHandler) Create(c *gin.Context) {
	var exam model.Exam
	if err := c.ShouldBindJSON(&exam); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.Create(&exam); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, exam)
}

// GetByID returns an exam by ID
func (h *ExamHandler) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid exam id"})
		return
	}
	exam, err := h.service.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "exam not found"})
		return
	}
	c.JSON(http.StatusOK, exam)
}

// Update updates an existing exam
func (h *ExamHandler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid exam id"})
		return
	}
	var exam model.Exam
	if err := c.ShouldBindJSON(&exam); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	exam.ID = uint(id)
	if err := h.service.Update(&exam); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, exam)
}

// Delete deletes an exam by ID
func (h *ExamHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid exam id"})
		return
	}
	if err := h.service.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

// List returns all exams
func (h *ExamHandler) List(c *gin.Context) {
	exams, err := h.service.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, exams)
}
