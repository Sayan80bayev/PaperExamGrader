package router

import (
	"PaperExamGrader/internal/config"
	"PaperExamGrader/internal/delivery"
	"PaperExamGrader/internal/middleware"
	"PaperExamGrader/internal/repository"
	"PaperExamGrader/internal/service"
	"PaperExamGrader/internal/storage"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupAnswerRoutes(r *gin.Engine, db *gorm.DB, cfg *config.Config, minio *storage.MinioStorage) {
	answerRepository := repository.GetAnswerRepository(db)
	answerService := service.NewAnswerService(answerRepository, minio)
	handler := delivery.NewAnswerHandler(answerService)

	answerGroup := r.Group("/answers", middleware.AuthMiddleware(cfg.JWTSecret))
	{
		answerGroup.POST("/upload", handler.Upload)
		answerGroup.POST("/upload_zip", handler.UploadFromZip)
		answerGroup.GET("/:id", handler.GetByID)
		answerGroup.GET("/exam/:exam_id", handler.GetByExamID)
		answerGroup.PUT("/:id/grade", handler.UpdateGrade)
		answerGroup.DELETE("/:id", handler.Delete)
	}
}
