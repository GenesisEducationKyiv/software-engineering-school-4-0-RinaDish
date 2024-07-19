package repo

import (
	"context"

	"github.com/RinaDish/currency-rates/internal/services"
)

type DBEmail struct {
    Email string
	IsActive bool
}

func (repo *Repository) SetEmail(ctx context.Context, email string) error {
	e := &DBEmail{
		Email: email,
		IsActive: true,
	}

	return repo.DB.Table("emails").Create(e).Error
}

func (repo *Repository) GetEmails(ctx context.Context) ([]services.Email, error) {
	result := make([]services.Email, 0)
	err := repo.DB.Table("emails").Where("is_active = ?", true).Find(&result).Error
	
	return result, err
}

func (repo *Repository) DeactivateEmail(ctx context.Context, email string) error {
	err := repo.DB.Model(services.Email{}).Where("email = ?", email).Update("is_active", false).Error

	if err != nil {
		return err
	}

	return nil
}
