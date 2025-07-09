// FILEPATH: /Users/sayanseksenbaev/Programming/PaperExamGrader/internal/router/exam_routes.go

package router

import (
	"PaperExamGrader/internal/config"
	"PaperExamGrader/internal/delivery"
	"PaperExamGrader/internal/middleware"
	"PaperExamGrader/internal/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupExamRoutes(r *gin.Engine, db *gorm.DB, cfg *config.Config) {
	examService := service.NewExamService(db)
	examHandler := delivery.NewExamHandler(examService)

	examGroup := r.Group("/exams", middleware.AuthMiddleware(cfg.JWTSecret))
	{
		examGroup.POST("/", examHandler.Create)
		examGroup.GET("/:id", examHandler.GetByID)
		examGroup.PUT("/:id", examHandler.Update)
		examGroup.DELETE("/:id", examHandler.Delete)
		examGroup.GET("/", examHandler.List)
	}
}
