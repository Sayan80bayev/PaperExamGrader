package bootstrap

import (
	"PaperExamGrader/internal/config"
	"PaperExamGrader/pkg/logging"
	"database/sql"
	"errors"
	"github.com/golang-migrate/migrate/v4"
	migrateps "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Container struct {
	DB *gorm.DB
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
		DB: db,
		//Redis:        redisClient,
		//Minio:        minioClient,
		//Producer:     producer,
		//Consumer:     consumer,
		Config: cfg,
		//Repositories: repositories,
	}, nil
}

func initDatabase(cfg *config.Config) (*gorm.DB, error) {
	logger := logging.GetLogger()

	db, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{})
	if err != nil {
		logger.Fatal("Error connecting to the database:", err)
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		logger.Fatal("Error getting generic DB object:", err)
		return nil, err
	}

	g, err2, done := initMigrations(sqlDB, logger)
	if done {
		return g, err2
	}

	return db, nil
}

func initMigrations(sqlDB *sql.DB, logger *logrus.Logger) (*gorm.DB, error, bool) {
	driver, err := migrateps.WithInstance(sqlDB, &migrateps.Config{})
	if err != nil {
		logger.Fatal("Error initializing migration driver:", err)
		return nil, err, true
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres", driver)
	if err != nil {
		logger.Fatal("Error creating migration instance:", err)
		return nil, err, true
	}

	version, _, err := m.Version()
	if err != nil {
		logger.Fatal("Error checking migration version:", err)
		return nil, err, true
	}

	if version == 0 {
		logger.Info("⚡ Applying migrations for the first time")
	} else {
		logger.Info("⚡ Database already migrated (version:", version, ")")
	}

	if version == 1 {
		logger.Warn("⚠️ Database is in dirty state, forcing migration to version 1")
		if err := m.Force(1); err != nil {
			logger.Fatal("Error forcing migration to version 1:", err)
			return nil, err, true
		}
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		logger.Fatal("Migration failed:", err)
		return nil, err, true
	}
	logger.Info("✅ Migrations applied successfully")
	return nil, nil, false
}
