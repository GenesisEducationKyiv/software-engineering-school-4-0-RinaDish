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

func (natsBroker *NATSBroker) Subscribe(subj string, cb nats.MsgHandler) (*nats.Subscription, error) {
	return natsBroker.conn.Subscribe(subj, cb)
}

func (natsBroker *NATSBroker) Drain() error {
	return natsBroker.conn.Drain()
}
