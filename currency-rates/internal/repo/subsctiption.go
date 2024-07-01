package repo

import (
	"context"

	"github.com/RinaDish/currency-rates/internal/services"
)

func (repo *Repository) SetEmail(ctx context.Context, email string) error {
	e := &services.Email{
		Email: email,
	}

	return repo.DB.Table("emails").Create(e).Error
}

func (repo *Repository) GetEmails(ctx context.Context) ([]services.Email, error) {
	result := make([]services.Email, 0)
	err := repo.DB.Table("emails").Find(&result).Error
	
	return result, err
}
