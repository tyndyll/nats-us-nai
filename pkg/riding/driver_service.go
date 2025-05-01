package riding

import (
	"context"
	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"log"
)

type DriverService struct {
	ID string
	nc *nats.Conn

	AcceptableLocations []string
}

func (ds *DriverService) Serve(ctx context.Context) {
	log.Println("driver service: starting for driver:", ds.ID)
	_, _ = ds.nc.Subscribe(DriverRideAcceptedTopic(ds.ID), ds.HandleRideAccepted)
	_, _ = ds.nc.Subscribe(DriverNotificationsTopic, ds.HandleRideRequest)

	// Subscribe to the ride topic
	<-ctx.Done()
}

func (ds *DriverService) HandleRideRequest(msg *nats.Msg) {
	log.Printf("driver service [%s]: Received a ride request message: %s", ds.ID, string(msg.Data))
	rr := &RideRequest{}
	if err := DecodeMessage(msg.Data, rr); err != nil {
		log.Printf("Error decoding Riding request: %v", err)
		return
	}

	// Acceptable location check

	// Register Interest
	interestMessage, err := EncodeMessage(&DriverInterested{
		ID:      ds.ID,
		Request: rr,
	})
	if err != nil {
		log.Printf("Error encoding driver interest message: %v", err)
		return
	}

	log.Printf("driver service [%s]: publishing driver interest message %s", ds.ID, string(interestMessage))
	if err = ds.nc.Publish(DriverInterestedTopic, interestMessage); err != nil {
		log.Printf("Error publishing driver interest: %v", err)
	}
}

func (ds *DriverService) HandleRideAccepted(msg *nats.Msg) {
	// We know this message is just a byte string

	rideID := string(msg.Data)
	log.Printf("driver service [%s]: Ride accepted: %v\n", ds.ID, rideID)
}

func NewDriverService(nc *nats.Conn) *DriverService {
	return &DriverService{
		ID: uuid.New().String(),
		nc: nc,
	}
}
