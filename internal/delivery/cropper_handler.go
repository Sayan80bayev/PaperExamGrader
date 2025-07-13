package delivery

import (
	"PaperExamGrader/internal/service"
	"PaperExamGrader/internal/transport/request"
	"fmt"
	"github.com/gin-gonic/gin"
)

type CropperHandler struct {
	manual *service.ManualCropper
}

func NewCropperHandler(manual *service.ManualCropper) *CropperHandler {
	return &CropperHandler{manual: manual}
}

func (h *CropperHandler) CropManual(c *gin.Context) {
	// 1. Parse examID from path
	examIDStr := c.Param("exam_id")
	if examIDStr == "" {
		c.JSON(400, gin.H{"error": "exam_id is required"})
		return
	}
	var examID uint
	if _, err := fmt.Sscanf(examIDStr, "%d", &examID); err != nil {
		c.JSON(400, gin.H{"error": "invalid exam_id"})
		return
	}

	// 2. Bind array of BBoxMetaDB from JSON body
	var dbBBoxes []request.BBoxMetaDB
	if err := c.ShouldBindJSON(&dbBBoxes); err != nil {
		c.JSON(400, gin.H{"error": "invalid JSON input", "details": err.Error()})
		return
	}

	// 3. Convert BBoxMetaDB -> BBoxMeta
	var bboxes []service.BBoxMeta
	for _, dbMeta := range dbBBoxes {
		meta, err := service.FromBBoxRequest(dbMeta)
		if err != nil {
			c.JSON(400, gin.H{
				"error":   "invalid bbox_percent format",
				"details": err.Error(),
			})
			return
		}
		bboxes = append(bboxes, meta)
	}

	// 4. Set crop parameters
	h.manual.SetCropParams(bboxes)

	// 5. Perform cropping
	urls, err := h.manual.CropFromExamPDFURLs(examID)
	if err != nil {
		c.JSON(500, gin.H{"error": "cropping failed", "details": err.Error()})
		return
	}

	// 6. Return result
	c.JSON(200, gin.H{
		"message":      "Cropping and upload successful",
		"image_urls":   urls,
		"croppedCount": len(urls),
	})
}
