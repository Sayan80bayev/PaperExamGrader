package router

import (
	"PaperExamGrader/internal/config"
	"PaperExamGrader/internal/delivery"
	"PaperExamGrader/internal/middleware"
	"PaperExamGrader/internal/repository"
	"PaperExamGrader/internal/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupBBoxRoutes(r *gin.Engine, db *gorm.DB, cfg *config.Config) {
	bboxRepo := repository.GetBBoxRepository(db)
	examRepo := repository.GetExamRepository(db)
	bboxService := service.NewBBoxService(bboxRepo, examRepo)
	bboxHandler := delivery.NewBBoxHandler(bboxService)

	api := r.Group("/api/bboxes", middleware.AuthMiddleware(cfg.JWTSecret))
	{
		api.POST("", bboxHandler.Create)
		api.GET("/list/:id", bboxHandler.GetAllByExamID)
		api.GET("/:id", bboxHandler.GetByID)
		api.PUT("/:id", bboxHandler.Update)
		api.DELETE("/:id", bboxHandler.Delete)
	}
}
