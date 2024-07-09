package queue

import (
	"github.com/nats-io/nats.go"

	"github.com/RinaDish/subscription-sender/tools"
)

type Queue struct {
	QueueConn *nats.Conn
	logger tools.Logger
}


func NewQueue(nats *nats.Conn, logger tools.Logger) (*Queue) {
	return &Queue{
		QueueConn:     nats,
		logger: logger.With("service", "queue"),
	}
}
