package repo

import (
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Repository struct {
	DB     *gorm.DB
	logger *zap.SugaredLogger
}

func NewAdminRepository(url string, logger *zap.SugaredLogger) (*Repository, error) {
	db, err := gorm.Open(postgres.Open(url), &gorm.Config{})
	
	return &Repository{
		DB:     db,
		logger: logger.With("service", "repository"),
	}, err
}
