package repo

import (
	"gorm.io/gorm"

	"github.com/RinaDish/subscription-sender/tools"
)

type Repository struct {
	DB     *gorm.DB
	logger tools.Logger
}

func NewAdminRepository(db *gorm.DB, logger tools.Logger) (*Repository) {
	return &Repository{
		DB:     db,
		logger: logger.With("service", "repository"),
	}
}
