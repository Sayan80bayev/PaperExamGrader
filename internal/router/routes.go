package router

import (
	"PaperExamGrader/internal/config"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(r *gin.Engine, db *gorm.DB, cfg *config.Config) {
	SetupManualCutterRoutes(r)
	SetupExamRoutes(r, db, cfg)
}
