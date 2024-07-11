package repo

import (
	"context"
	"time"

	"github.com/RinaDish/subscription-sender/internal/services"
)

func (repo *Repository) SetMessages(ctx context.Context, message services.Message) error {
	return repo.DB.Table("messages").Create(&message).Error
}

func (repo *Repository) GetMessages(ctx context.Context) ([]services.Message, error) {
	result := make([]services.Message, 0)

	currentTime := time.Now()

	err := repo.DB.Table("messages").Where("sent = ?", false).Where("sending_time < ?", currentTime).Find(&result).Error
	
	return result, err
}

func (repo *Repository) UpdateMessages(ctx context.Context, message services.Message) error {
	err := repo.DB.Model(services.Message{}).Where("id = ?", message.ID).Update("sent", true).Error

	if err != nil {
		return err
	}

	return nil
}



