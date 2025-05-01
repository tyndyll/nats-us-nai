package riding

import "github.com/nats-io/nats.go"

func NewNATSConnection(srv string) (*nats.Conn, error) {
	if srv == "" {
		srv = nats.DefaultURL
	}
	nc, err := nats.Connect(srv)
	if err != nil {
		return nil, err
	}
	return nc, nil
}
