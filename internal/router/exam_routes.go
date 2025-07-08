// FILEPATH: /Users/sayanseksenbaev/Programming/PaperExamGrader/internal/router/exam_routes.go

package router

import (
	"PaperExamGrader/internal/delivery"
	"PaperExamGrader/internal/service"
	"github.com/gin-gonic/gin"
)

func SetupExamRoutes(r *gin.Engine, examService *service.AnswerService) {
	examHandler := delivery.NewExamHandler(examService)

	examGroup := r.Group("/exams")
	{
		examGroup.POST("/", examHandler.Create)
		examGroup.GET("/:id", examHandler.GetByID)
		examGroup.PUT("/:id", examHandler.Update)
		examGroup.DELETE("/:id", examHandler.Delete)
		examGroup.GET("/", examHandler.List)
	}
}
