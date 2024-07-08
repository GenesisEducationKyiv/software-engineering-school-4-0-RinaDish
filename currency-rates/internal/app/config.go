package app

type Config struct {
	Address string `envconfig:"ADDRESS"`
	NotificationServiceURL string `envconfig:"NOTIFICATION_SERVICE_URL"`
	DBURL string `envconfig:"DB_URL"`
	NatsURL string `envconfig:"NATS_URL"`
}
