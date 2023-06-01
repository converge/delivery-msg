package handlers

import (
	"context"
	"delivery-msg/internal/domain"
	"delivery-msg/internal/repositories"
	"delivery-msg/internal/services"
	"delivery-msg/pb"
	"github.com/charmbracelet/log"
	"github.com/google/uuid"
	"time"
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
	now := time.Now().Format("2006-01-02 15:04:05")
	if err := deliveryHandler.dbRepository.CreateDelivery(in.SourceAddress, in.DestinationAddress, trackingCode, now, now); err != nil {
		return nil, err
	}

	err := deliveryHandler.natsClient.Publish(
		trackingCode,
		in.SourceAddress,
		in.DestinationAddress,
		pb.StatusEnum_CONFIRMED,
		now,
		now,
	)
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

	log.Info(in.TrackingCode)

	modified := time.Now().Format("2006-01-02 15:04:05")
	updatedDelivery, err := deliveryHandler.dbRepository.UpdateDelivery(in.TrackingCode, modified, in.Status.String())
	if err != nil {
		log.Error(err)
		return nil, err
	}

	// todo: centralize format
	now := time.Now().Format("2006-01-02 15:04:05")

	err = deliveryHandler.natsClient.Publish(
		in.TrackingCode,
		updatedDelivery.SourceAddress,
		updatedDelivery.DestinationAddress,
		in.GetStatus(),
		updatedDelivery.Created,
		now,
	)
	if err != nil {
		return nil, err
	}

	return &pb.UpdateDeliveryResponse{}, nil
}

func (deliveryHandler *DeliveryHandler) GetDelivery() (*[]domain.Delivery, error) {

	dbData, err := deliveryHandler.dbRepository.GetAllDeliveries()
	if err != nil {
		return nil, err
	}

	return &dbData, nil
}
