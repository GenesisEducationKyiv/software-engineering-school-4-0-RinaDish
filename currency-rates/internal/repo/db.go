package repo

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Repository struct {
	DB     *gorm.DB
	logger *zap.SugaredLogger
}

func NewAdminRepository(db *gorm.DB, l *zap.SugaredLogger) (*Repository) {
	return &Repository{
		DB:     db,
		logger: l.With("service", "repository"),
	}
}
