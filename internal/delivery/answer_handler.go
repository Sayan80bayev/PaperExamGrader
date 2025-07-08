// FILEPATH: /Users/sayanseksenbaev/Programming/PaperExamGrader/internal/handler/exam_handler.go

package delivery

import (
	"net/http"
	"strconv"

	"PaperExamGrader/internal/model"
	"PaperExamGrader/internal/service"
	"github.com/gin-gonic/gin"
)

type AnswerHandler struct {
	service *service.AnswerService
}

func NewAnswerHandler(service *service.AnswerService) *AnswerHandler {
	return &AnswerHandler{service: service}
}

// Create creates a new Answer
func (h *AnswerHandler) Create(c *gin.Context) {
	var Answer model.Answer
	if err := c.ShouldBindJSON(&Answer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.Create(&Answer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, Answer)
}

// GetByID returns an answer by ID
func (h *AnswerHandler) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid answer id"})
		return
	}
	answer, err := h.service.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "answer not found"})
		return
	}
	c.JSON(http.StatusOK, answer)
}

// Update updates an existing answer
func (h *AnswerHandler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid answer id"})
		return
	}
	var answer model.Answer
	if err := c.ShouldBindJSON(&answer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	answer.ID = uint(id)
	if err := h.service.Update(&answer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, answer)
}

// Delete deletes an answer by ID
func (h *AnswerHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid answer id"})
		return
	}
	if err := h.service.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

// List returns all answers
func (h *AnswerHandler) List(c *gin.Context) {
	answers, err := h.service.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, answers)
}
