package app

type Config struct {
	Address string `envconfig:"ADDRESS"`
	DBUrl string `envconfig:"DB_URL"`
	EmailAddress string `envconfig:"EMAIL_ADDRESS"`
	EmailPass string `envconfig:"EMAIL_PASS"`
}
