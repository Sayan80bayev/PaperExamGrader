// FILEPATH: /Users/sayanseksenbaev/Programming/PaperExamGrader/internal/router/class_routes.go

package router

import (
	"PaperExamGrader/internal/delivery"
	"PaperExamGrader/internal/service"
	"github.com/gin-gonic/gin"
)

func SetupClassRoutes(r *gin.Engine, classService *service.ExamService) {
	classHandler := delivery.NewClassHandler(classService)

	classGroup := r.Group("/classes")
	{
		classGroup.POST("/", classHandler.Create)
		classGroup.GET("/:id", classHandler.GetByID)
		classGroup.PUT("/:id", classHandler.Update)
		classGroup.DELETE("/:id", classHandler.Delete)
		classGroup.GET("/", classHandler.List)
	}
}
