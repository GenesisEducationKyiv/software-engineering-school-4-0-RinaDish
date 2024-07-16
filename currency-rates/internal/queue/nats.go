package queue

import (
	"github.com/nats-io/nats.go"
)

type NATSBroker struct {
	conn *nats.Conn
}

func NewNATSBroker(conn *nats.Conn) *NATSBroker {
	return &NATSBroker{
		conn: conn,
	}
}

func (natsBroker *NATSBroker) Publish(subj string, data []byte) error {
	return natsBroker.conn.Publish(subj, data)
}

func (natsBroker *NATSBroker) Drain() error {
	return natsBroker.conn.Drain()
}
