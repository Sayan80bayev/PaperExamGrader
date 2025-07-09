// FILEPATH: /Users/sayanseksenbaev/Programming/PaperanswersGrader/internal/router/exam_routes.go

package router

import (
	"PaperExamGrader/internal/delivery"
	"PaperExamGrader/internal/service"
	"github.com/gin-gonic/gin"
)

func SetupAnswerRoutes(r *gin.Engine, answerService *service.AnswerService) {
	answerHandler := delivery.NewAnswerHandler(answerService)
	
	answerGroup := r.Group("/answers")
	{
		answerGroup.POST("/", answerHandler.Create)
		answerGroup.GET("/:id", answerHandler.GetByID)
		answerGroup.PUT("/:id", answerHandler.Update)
		answerGroup.DELETE("/:id", answerHandler.Delete)
		answerGroup.GET("/", answerHandler.List)
	}
}
