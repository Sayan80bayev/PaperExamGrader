package delivery

import (
	"PaperExamGrader/internal/model"
	"PaperExamGrader/internal/service"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type BBoxHandler struct {
	service service.BBoxService
}

func NewBBoxHandler(service service.BBoxService) *BBoxHandler {
	return &BBoxHandler{service: service}
}

func (h *BBoxHandler) Create(c *gin.Context) {
	var input model.BBoxMetaDB
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate bbox size
	var bboxArr []float64
	if err := json.Unmarshal(input.BBoxPercent, &bboxArr); err != nil || len(bboxArr) != 4 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "BBoxPercent must be a JSON array of 4 float64 values"})
		return
	}

	if err := h.service.Create(&input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, input)
}

func (h *BBoxHandler) GetAll(c *gin.Context) {
	entries, err := h.service.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, entries)
}

func (h *BBoxHandler) GetByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	entry, err := h.service.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
		return
	}
	c.JSON(http.StatusOK, entry)
}

func (h *BBoxHandler) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var input model.BBoxMetaDB
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	input.ID = uint(id)
	if err := h.service.Update(&input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, input)
}

func (h *BBoxHandler) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := h.service.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
