package app

type Config struct {
	Address string `envconfig:"ADDRESS"`
	NotificationServiceUrl string `envconfig:"NOTIFICATION_SERVISE_URL"`
	DBUrl string `envconfig:"DB_URL"`
}
