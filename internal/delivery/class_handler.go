// FILEPATH: /Users/sayanseksenbaev/Programming/PaperExamGrader/internal/handler/class_handler.go

package delivery

import (
	"net/http"
	"strconv"

	"PaperExamGrader/internal/model"
	"PaperExamGrader/internal/service"
	"github.com/gin-gonic/gin"
)

type ClassHandler struct {
	service *service.ClassService
}

func NewClassHandler(service *service.ClassService) *ClassHandler {
	return &ClassHandler{service: service}
}

// Create creates a new class
func (h *ClassHandler) Create(c *gin.Context) {
	var class model.Class
	if err := c.ShouldBindJSON(&class); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.Create(&class); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, class)
}

// GetByID returns a class by ID
func (h *ClassHandler) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid class id"})
		return
	}
	class, err := h.service.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "class not found"})
		return
	}
	c.JSON(http.StatusOK, class)
}

// Update updates an existing class
func (h *ClassHandler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid class id"})
		return
	}
	var class model.Class
	if err := c.ShouldBindJSON(&class); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	class.ID = uint(id)
	if err := h.service.Update(&class); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, class)
}

// Delete deletes a class by ID
func (h *ClassHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid class id"})
		return
	}
	if err := h.service.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

// List returns all classes
func (h *ClassHandler) List(c *gin.Context) {
	classes, err := h.service.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, classes)
}
