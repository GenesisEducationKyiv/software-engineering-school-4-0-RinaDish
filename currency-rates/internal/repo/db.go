package repo

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/RinaDish/currency-rates/tools"
)

type Repository struct {
	DB     *gorm.DB
	logger logger.Logger
}

func NewAdminRepository(url string, logger logger.Logger) (*Repository, error) {
	db, err := gorm.Open(postgres.Open(url), &gorm.Config{})
	
	return &Repository{
		DB:     db,
		logger: logger.With("service", "repository"),
	}, err
}
