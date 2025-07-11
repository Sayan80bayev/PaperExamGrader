package router

import (
	"PaperExamGrader/internal/delivery"
	"PaperExamGrader/internal/repository"
	"PaperExamGrader/internal/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupBBoxRoutes(r *gin.Engine, db *gorm.DB) {
	bboxRepo := repository.GetBBoxRepository(db)
	bboxService := service.NewBBoxService(bboxRepo)
	bboxHandler := delivery.NewBBoxHandler(bboxService)

	api := r.Group("/api/bboxes")
	{
		api.POST("", bboxHandler.Create)
		api.GET("", bboxHandler.GetAll)
		api.GET("/:id", bboxHandler.GetByID)
		api.PUT("/:id", bboxHandler.Update)
		api.DELETE("/:id", bboxHandler.Delete)
	}
}
