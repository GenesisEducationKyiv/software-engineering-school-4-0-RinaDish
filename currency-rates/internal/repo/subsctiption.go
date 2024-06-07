package repo

import "context"

type Email struct {
	ID    int    `json:"id" gorm:"id"`
	Email string `json:"email" gorm:"email"`
}

func (r *Repository) SetEmail(ctx context.Context, email string) error {
	e := &Email{
		Email: email,
	}
	return r.DB.Table("emails").Create(e).Error
}

func (r *Repository) GetEmails(ctx context.Context) ([]Email, error) {
	result := make([]Email, 0)
	err := r.DB.Table("emails").Find(&result).Error
	return result, err
}
