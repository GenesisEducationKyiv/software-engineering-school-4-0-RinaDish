package app

type Config struct {
	Address string `envconfig:"ADDRESS"`
	EmailAddress string `envconfig:"EMAIL_ADDRESS"`
	EmailPass string `envconfig:"EMAIL_PASS"`
	NatsURL string `envconfig:"NATS_URL"`
}
