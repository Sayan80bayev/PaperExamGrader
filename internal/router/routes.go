package router

import (
	"PaperExamGrader/internal/config"
	"PaperExamGrader/internal/storage"
	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"gorm.io/gorm"
)

func SetupRoutes(r *gin.Engine, db *gorm.DB, cfg *config.Config, minioClient *minio.Client) {
	minio := &storage.MinioStorage{
		Client: minioClient,
		Cfg:    cfg,
	}

	SetupBBoxRoutes(r, db, cfg)
	SetupManualCutterRoutes(r, db, minio)
	SetupExamRoutes(r, db, cfg)
	SetupAnswerRoutes(r, db, cfg, minio)
}
