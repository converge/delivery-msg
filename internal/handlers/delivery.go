package handlers

import (
	"context"
	"delivery-msg/internal/repositories"
	"delivery-msg/internal/services"
	"delivery-msg/pb"
	"github.com/google/uuid"
)

type DeliveryHandler struct {
	pb.UnimplementedDeliveryServiceServer
	dbRepository *repositories.PostgreSQLRepository
	natsClient   *services.NATSClient
}

func NewDeliveryHandler(
	dbRepository *repositories.PostgreSQLRepository,
	natsClient *services.NATSClient,
) *DeliveryHandler {

	return &DeliveryHandler{
		dbRepository: dbRepository,
		natsClient:   natsClient,
	}
}

func (deliveryHandler *DeliveryHandler) CreateDelivery(
	ctx context.Context,
	in *pb.CreateDeliveryRequest,
) (*pb.CreateDeliveryResponse, error) {

	trackingCode := uuid.NewString()
	if err := deliveryHandler.dbRepository.CreateDelivery(in.SourceAddress, in.DestinationAddress, trackingCode); err != nil {
		return nil, err
	}

	err := deliveryHandler.natsClient.Publish(trackingCode, pb.StatusEnum_CONFIRMED)
	if err != nil {
		return nil, err
	}

	return &pb.CreateDeliveryResponse{
		TrackingCode:       trackingCode,
		SourceAddress:      in.SourceAddress,
		DestinationAddress: in.DestinationAddress,
		Status:             pb.StatusEnum_CONFIRMED,
	}, nil
}

func (deliveryHandler *DeliveryHandler) UpdateDelivery(
	ctx context.Context,
	in *pb.UpdateDeliveryRequest,
) (*pb.UpdateDeliveryResponse, error) {

	if err := deliveryHandler.dbRepository.UpdateDelivery(in.TrackingCode, in.Status.String()); err != nil {
		return nil, err
	}

	err := deliveryHandler.natsClient.Publish(in.TrackingCode, in.GetStatus())
	if err != nil {
		return nil, err
	}

	return &pb.UpdateDeliveryResponse{}, nil
}
