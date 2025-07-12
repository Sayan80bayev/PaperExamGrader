package delivery

import (
	"PaperExamGrader/internal/service"
	"encoding/json"
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
	// 1. Parse `answerID` from form
	answerIDStr := c.PostForm("answer_id")
	if answerIDStr == "" {
		c.JSON(400, gin.H{"error": "answer_id is required"})
		return
	}
	var answerID uint
	_, err := fmt.Sscanf(answerIDStr, "%d", &answerID)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid answer_id"})
		return
	}

	// 2. Parse `bbox_params` JSON from form
	bboxJSON := c.PostForm("bbox_params")
	if bboxJSON == "" {
		c.JSON(400, gin.H{"error": "bbox_params is required"})
		return
	}

	var bboxes []service.BBoxMeta
	if err := json.Unmarshal([]byte(bboxJSON), &bboxes); err != nil {
		c.JSON(400, gin.H{"error": "invalid bbox_params format", "details": err.Error()})
		return
	}

	// 3. Parse uploaded files
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(400, gin.H{"error": "failed to parse multipart form"})
		return
	}
	files := form.File["files"]
	if len(files) == 0 {
		c.JSON(400, gin.H{"error": "no files uploaded"})
		return
	}

	// 4. Set crop parameters
	h.manual.SetCropParams(bboxes)

	// 5. Perform cropping
	urls, err := h.manual.CropFromMultipartPDFs(files, answerID)
	if err != nil {
		c.JSON(500, gin.H{"error": "cropping failed", "details": err.Error()})
		return
	}

	// 6. Return success with uploaded image URLs
	c.JSON(200, gin.H{
		"message":      "Cropping and upload successful",
		"image_urls":   urls,
		"croppedCount": len(urls),
	})
}
