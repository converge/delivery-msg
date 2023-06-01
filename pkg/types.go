package pkg

import (
	"delivery-msg/internal/domain"
	"github.com/nats-io/nats.go"
)

type DeliveryData []domain.Delivery
type NatsListener *nats.Msg
