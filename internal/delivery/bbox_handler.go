package delivery

import (
	"PaperExamGrader/internal/service"
	"PaperExamGrader/internal/transport/request"
	"net/http"
	"strconv"
	
	"github.com/gin-gonic/gin"
)

type BBoxHandler struct {
	service *service.BBoxService
}

func NewBBoxHandler(service *service.BBoxService) *BBoxHandler {
	return &BBoxHandler{service: service}
}

func (h *BBoxHandler) CreateTemplate(c *gin.Context) {
	var req request.CreateBBoxTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	instructorId := c.GetUint("user_id")
	
	templateResponse, err := h.service.CreateTemplateWithBBoxes(req, instructorId)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusCreated, templateResponse)
}

func (h *BBoxHandler) GetTemplatesByExamID(c *gin.Context) {
	examID, _ := strconv.Atoi(c.Param("id"))
	instructorId := c.GetUint("user_id")
	
	templates, err := h.service.GetTemplatesByExamID(uint(examID), instructorId)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, templates)
}

func (h *BBoxHandler) DeleteTemplate(c *gin.Context) {
	templateID, _ := strconv.Atoi(c.Param("id"))
	instructorId := c.GetUint("user_id")
	
	if err := h.service.DeleteTemplate(uint(templateID), instructorId); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}
	
	c.Status(http.StatusNoContent)
}
