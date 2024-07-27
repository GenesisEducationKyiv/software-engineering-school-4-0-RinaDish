package repo

import (
	"context"
	"github.com/lib/pq"
	"time"

	"github.com/RinaDish/subscription-sender/internal/services"
)

type message struct {
	ID int `gorm:"id"`
	Rate float64 `gorm:"rate"`
	Emails pq.StringArray `gorm:"type:text[]"`
	CreatedAt time.Time `gorm:"created_at"`
	EventID uint8 `gorm:"event_id"`
	EventType string `gorm:"event_type"`
	SendingTime time.Time `gorm:"sending_time"`
	Sent bool `gorm:"sent"`
}

func (repo *Repository) SaveMessages(ctx context.Context, m services.Message) error {
	msg := &message{
		ID: m.ID,
		Rate: m.Rate,
		Emails: pq.StringArray(m.Emails),
		CreatedAt: m.CreatedAt,
		EventID: m.EventID,
		EventType: m.EventType,
		SendingTime: m.SendingTime,
		Sent: m.Sent,
	}

	return repo.DB.Table("messages").Create(&msg).WithContext(ctx).Error
}

func (repo *Repository) GetReadyToSendMessages(ctx context.Context) ([]services.Message, error) {
	msgsList := make([]message, 0)

	currentTime := time.Now()

	err := repo.DB.Table("messages").WithContext(ctx).Where("sent = ?", false).Where("sending_time < ?", currentTime).Find(&msgsList).Error
	
	if err != nil {
		return nil, err
	}

	result := make([]services.Message, 0, len(msgsList))
	for _, msg := range msgsList {
		result = append(result, services.Message{
			ID: msg.ID,
			Rate: msg.Rate,
			Emails: msg.Emails,
			CreatedAt: msg.CreatedAt,
			EventID: msg.EventID,
			EventType: msg.EventType,
			SendingTime: msg.SendingTime,
			Sent: msg.Sent,
		})
	}

	return result, err
}

func (repo *Repository) SetSendStatus(ctx context.Context, msgID int, send bool) error {
	return repo.DB.Model(&message{}).WithContext(ctx).Where("id = ?", msgID).Update("sent", send).Error
}
