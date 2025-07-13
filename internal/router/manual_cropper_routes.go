package router

import (
	"PaperExamGrader/internal/delivery"
	"PaperExamGrader/internal/repository"
	"PaperExamGrader/internal/service"
	"PaperExamGrader/internal/storage"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupManualCutterRoutes(r *gin.Engine, db *gorm.DB, minio *storage.MinioStorage) {
	imgRepo := repository.GetImageRepository(db)
	answerRepo := repository.GetAnswerRepository(db)
	mcService := service.NewManualCropper(minio, imgRepo, answerRepo)
	mcHandler := delivery.NewCropperHandler(mcService)

	r.POST("/api/v1/crop/manual/:exam_id", func(c *gin.Context) {
		mcHandler.CropManual(c)
	})
}
