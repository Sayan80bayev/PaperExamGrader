package bootstrap

import (
	"PaperExamGrader/internal/storage"
	"database/sql"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/sirupsen/logrus"

	migrateps "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/minio/minio-go/v7"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"PaperExamGrader/internal/config"
	"PaperExamGrader/pkg/logging"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type Container struct {
	DB     *gorm.DB
	Minio  *minio.Client
	Config *config.Config
}

func Init() (*Container, error) {
	logger := logging.GetLogger()
	fmt.Println("\n .oooooo..o       .o.       oooooo   oooo       .o.       ooooo      ooo \nd8P'    `Y8      .888.       `888.   .8'       .888.      `888b.     `8' \nY88bo.          .8\"888.       `888. .8'       .8\"888.      8 `88b.    8  \n `\"Y8888o.     .8' `888.       `888.8'       .8' `888.     8   `88b.  8  \n     `\"Y88b   .88ooo8888.       `888'       .88ooo8888.    8     `88b.8  \noo     .d8P  .8'     `888.       888       .8'     `888.   8       `888  \n8\"\"88888P'  o88o     o8888o     o888o     o88o     o8888o o8o        `8  \n                                                                         \n")
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Fatal("Error loading configuration:", err)
		return nil, err
	}

	db, err := initDatabase(cfg)
	if err != nil {
		return nil, err
	}

	minioClient := storage.Init(cfg)

	logger.Info("✅ Dependencies initialized successfully")

	return &Container{
		DB:     db,
		Minio:  minioClient,
		Config: cfg,
	}, nil
}

func initDatabase(cfg *config.Config) (*gorm.DB, error) {
	logger := logging.GetLogger()

	db, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{})
	if err != nil {
		logger.Fatal("Error connecting to the database:", err)
		return nil, err
	}

	//sqlDB, err := db.DB()
	//if err != nil {
	//	logger.Fatal("Error getting generic DB object:", err)
	//	return nil, err
	//}

	//g, err2, done := initMigrations(sqlDB, logger)
	//if done {
	//	return g, err2
	//}

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
