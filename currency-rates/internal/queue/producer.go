package queue

import (
	"github.com/nats-io/nats.go"

	"github.com/RinaDish/currency-rates/tools"
)

type Queue struct {
	Nats *nats.Conn
	logger tools.Logger
}


func NewQueue(nats *nats.Conn, logger tools.Logger) (*Queue) {
	return &Queue{
		Nats: nats,
		logger: logger.With("service", "queue"),
	}
}
