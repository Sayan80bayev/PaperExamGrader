// FILEPATH: /Users/sayanseksenbaev/Programming/PaperExamGrader/internal/router/user_routes.go

package router

import (
	"PaperExamGrader/internal/delivery"
	"PaperExamGrader/internal/service"
	"github.com/gin-gonic/gin"
)

func SetupUserRoutes(r *gin.Engine, userService *service.UserService) {
	userHandler := delivery.NewUserHandler(userService)

	userGroup := r.Group("/users")
	{
		userGroup.POST("/", userHandler.Create)
		userGroup.GET("/:id", userHandler.GetByID)
		userGroup.PUT("/:id", userHandler.Update)
		userGroup.DELETE("/:id", userHandler.Delete)
		userGroup.GET("/", userHandler.List)
	}
}
