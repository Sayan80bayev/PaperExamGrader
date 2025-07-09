package main

import (
	"PaperExamGrader/internal/bootstrap"
	"PaperExamGrader/internal/router"
	"PaperExamGrader/pkg/logging"
	"github.com/gin-gonic/gin"
)

// Initialize logrus as the main logger
var logger = logging.GetLogger()

func main() {
	bs, err := bootstrap.Init()
	if err != nil {
		logger.Errorf("bootstrap init err: %v", err)
	}

	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = logger.Out
	r := gin.New()

	r.Use(gin.Recovery())
	r.Use(logging.Middleware)
	router.SetupRoutes(r, bs.DB, bs.Config)

	logger.Infof("🚀 Server is running on port %s", bs.Config.Port)
	err = r.Run(":" + bs.Config.Port)
	if err != nil {
		logger.Errorf("Error starting server: %v", err)
	}
}
