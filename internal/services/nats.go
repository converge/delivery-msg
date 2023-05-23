package services

import (
	"delivery-msg/internal/domain"
	"delivery-msg/pb"
	"encoding/json"
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

func (natsClient *NATSClient) Publish(
	trackingCode string,
	sourceAddress string,
	destinationAddress string,
	status pb.StatusEnum,
	created string,
	modified string,
) error {

	deliveryData := domain.Delivery{
		TrackingCode:       trackingCode,
		SourceAddress:      sourceAddress,
		DestinationAddress: destinationAddress,
		Status:             status.String(),
		Created:            created,
		Modified:           modified,
	}
	jsonData, err := json.Marshal(deliveryData)
	if err != nil {
		return err
	}

	err = natsClient.Conn.Publish("nats_development", jsonData)
	if err != nil {
		return err
	}

	return nil
}
