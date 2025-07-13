package main

import (
	"PaperExamGrader/internal/bootstrap"
	"PaperExamGrader/internal/router"
	"PaperExamGrader/pkg/logging"
	"github.com/gin-contrib/cors"
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
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	corsConfig.ExposeHeaders = []string{"Content-Length"}
	corsConfig.AllowCredentials = true

	r.Use(cors.New(corsConfig))
	r.RedirectTrailingSlash = false
	r.Use(gin.Recovery())
	r.Use(logging.Middleware)
	router.SetupRoutes(r, bs.DB, bs.Config, bs.Minio)

	logger.Infof("🚀 Server is running on port %s", bs.Config.Port)
	err = r.Run(":" + bs.Config.Port)
	if err != nil {
		logger.Errorf("Error starting server: %v", err)
	}
}
