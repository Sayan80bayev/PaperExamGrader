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
	bboxTemplateRepo := repository.GetBBoxTemplateRepository(db)
	examRepo := repository.GetExamRepository(db)

	bboxService := service.NewBBoxService(db, bboxRepo, bboxTemplateRepo, examRepo)
	bboxHandler := delivery.NewBBoxHandler(bboxService)

	api := r.Group("/api/bboxes", middleware.AuthMiddleware(cfg.JWTSecret))
	{
		api.POST("/template", bboxHandler.CreateTemplate)
		api.GET("/template/list/:id", bboxHandler.GetTemplatesByExamID)
		api.DELETE("/template/:id", bboxHandler.DeleteTemplate)
	}
}
