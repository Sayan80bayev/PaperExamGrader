package router

import (
	"PaperExamGrader/internal/config"
	"PaperExamGrader/internal/storage"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(r *gin.Engine, db *gorm.DB, cfg *config.Config) {
	minioClient := storage.Init(cfg)
	minio := &storage.MinioStorage{
		Client: minioClient,
		Cfg:    cfg,
	}
	SetupBBoxRoutes(r, db)
	SetupManualCutterRoutes(r, db, minio)
	SetupExamRoutes(r, db, cfg)
	SetupAnswerRoutes(r, db, cfg, minio)
}
