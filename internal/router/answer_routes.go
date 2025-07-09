package router

import (
	"PaperExamGrader/internal/delivery"
	"PaperExamGrader/internal/service"
	"github.com/gin-gonic/gin"
)

func SetupAnswerRoutes(r *gin.Engine, answerService *service.AnswerService) {
	handler := delivery.NewAnswerHandler(answerService)

	answerGroup := r.Group("/answers")
	{
		answerGroup.POST("/upload", handler.Upload)
		answerGroup.GET("/:id", handler.GetByID)
		answerGroup.GET("/exam/:exam_id", handler.GetByExamID)
		answerGroup.PUT("/:id/grade", handler.UpdateGrade)
		answerGroup.DELETE("/:id", handler.Delete)
	}
}
