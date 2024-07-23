package repo

import (
	"gorm.io/gorm"

	"github.com/RinaDish/currency-rates/tools"
)

type Repository struct {
	DB     *gorm.DB
	metrics Metrics
	logger tools.Logger
}

type Metrics interface {
    MonitorDBQuery(queryName string, queryFunc func() error) error
}

func NewAdminRepository(db *gorm.DB, logger tools.Logger, metrics Metrics) (*Repository) {
	return &Repository{
		DB:     db,
		metrics: metrics,
		logger: logger.With("service", "repository"),
	}
}
