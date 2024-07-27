package repo

import (
	"context"

	"github.com/RinaDish/currency-rates/internal/services"
)

type DBEmail struct {
    Email string
	IsActive bool
}

func (repo *Repository) SetEmail(ctx context.Context, email string, isActive bool) error {
	e := &DBEmail{
		Email: email,
		IsActive: isActive,
	}

	return repo.DB.Table("emails").Save(e).Error
}

func (repo *Repository) GetAllActiveEmails(ctx context.Context) ([]services.Email, error) {
	result := make([]services.Email, 0)
	err := repo.DB.Table("emails").Where("is_active = ?", true).Find(&result).Error
	
	return result, err
}
