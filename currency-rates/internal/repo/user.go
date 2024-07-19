package repo

import (
	"context")

type DBUser struct {
    Email string
}


func (repo *Repository) SetUser(ctx context.Context, email string) error {
	u := &DBUser{
		Email: email,
	}

	return repo.DB.Table("users").Create(u).Error
}

func (repo *Repository) DeleteUser(ctx context.Context, email string) error {
	return repo.DB.Table("users").Where("email = ?", email).Delete(&DBUser{}).Error
}
