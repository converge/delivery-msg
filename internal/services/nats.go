package services

import (
	"delivery-msg/pb"
	"github.com/nats-io/nats.go"
)

type NATSClient struct {
	Conn *nats.Conn
}

func NewNATSClient(conn *nats.Conn) (NATSClient, error) {
	return NATSClient{
		Conn: conn,
	}, nil
}

func (natsClient *NATSClient) Publish(orderId string, status pb.StatusEnum) error {

	err := natsClient.Conn.Publish("nats_development", []byte("orderId: "+orderId+" status: "+status.String()))
	if err != nil {
		return err
	}

	return nil
}
