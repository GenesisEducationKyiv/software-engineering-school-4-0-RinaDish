package repo

import "context"

type Email struct {
	ID    int    `json:"id" gorm:"id"`
	Email string `json:"email" gorm:"email"`
}

func (repo *Repository) SetEmail(ctx context.Context, email string) error {
	e := &Email{
		Email: email,
	}

	return repo.DB.Table("emails").Create(e).Error
}

func (repo *Repository) GetEmails(ctx context.Context) ([]Email, error) {
	result := make([]Email, 0)
	err := repo.DB.Table("emails").Find(&result).Error
	
	return result, err
}
