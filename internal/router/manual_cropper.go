package router

import (
	"PaperExamGrader/internal/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func SetupManualCutterRoutes(r *gin.Engine) {

	r.POST("/api/v1/crop/manual", func(c *gin.Context) {
		form, err := c.MultipartForm()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid multipart form"})
			return
		}

		files := form.File["pdfs"]
		if len(files) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No PDFs provided"})
			return
		}

		// Define permanent batch and output directories
		basePath := "/Users/sayanseksenbaev/Programming/PaperExamGrader"
		batchDir := filepath.Join(basePath, "temp")
		outputDir := filepath.Join(basePath, "answers")

		// Ensure directories exist
		if err := os.MkdirAll(batchDir, os.ModePerm); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create batch directory"})
			return
		}
		if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create output directory"})
			return
		}

		// Save uploaded PDF files
		for _, file := range files {
			dst := filepath.Join(batchDir, filepath.Base(file.Filename))
			if err := c.SaveUploadedFile(file, dst); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to save %s", file.Filename)})
				return
			}
		}

		// Initialize the manual cropper
		cropper := &service.ManualCropper{}
		cropper.SetCropParams([]service.BBoxMeta{
			{Page: 0, BBoxPercent: [4]float64{0, 0.12, 1, 0.32}},
			{Page: 0, BBoxPercent: [4]float64{0, 0.25, 1, 0.43}},
			{Page: 0, BBoxPercent: [4]float64{0, 0.37, 1, 0.57}},
			{Page: 0, BBoxPercent: [4]float64{0, 0.47, 1, 0.73}},
			{Page: 0, BBoxPercent: [4]float64{0, 0.65, 1, 1}},
		})

		// Run cropping with the permanent directories
		if err := cropper.Crop(batchDir, outputDir); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Collect result image paths
		entries, err := os.ReadDir(outputDir)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read output directory"})
			return
		}

		var results []string
		for _, entry := range entries {
			if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".jpg") {
				results = append(results, entry.Name())
			}
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Batch cropping completed.",
			"outputs": results,
		})
	})

}
