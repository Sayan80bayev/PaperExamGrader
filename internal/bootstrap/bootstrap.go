package bootstrap

import (
	"PaperExamGrader/internal/config"
	"PaperExamGrader/pkg/logging"
)

type Container struct {
	//DB           *gorm.DB
	//Redis        *redis.Client
	//Minio        *minio.Client
	//Producer     messaging.Producer
	//Consumer     messaging.Consumer
	Config       *config.Config
	Repositories map[string]interface{}
}

func Init() (*Container, error) {
	logger := logging.GetLogger()

	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Fatal("Error loading configuration:", err)
		return nil, err
	}

	return &Container{
		//DB:           db,
		//Redis:        redisClient,
		//Minio:        minioClient,
		//Producer:     producer,
		//Consumer:     consumer,
		Config: cfg,
		//Repositories: repositories,
	}, nil
}
