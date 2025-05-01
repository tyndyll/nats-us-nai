package riding

import (
	"context"
	"github.com/nats-io/nats.go"
	"log"
)

type RideRequestService struct {
	nc *nats.Conn

	sub *nats.Subscription
}

func (rrs *RideRequestService) Serve(ctx context.Context) error {
	log.Println("ride request service: starting")
	sub, err := rrs.nc.Subscribe(RideRequestTopic, func(msg *nats.Msg) {
		log.Println("ride request service: received ride request", string(msg.Data))
		// Advertise the ride request to drivers
		log.Println("ride request service: publishing ride request to all drivers:", string(msg.Data))
		if err := rrs.nc.Publish(DriverNotificationsTopic, msg.Data); err != nil {
			log.Printf("ride request service: error publishing ride request: %v", err)
		}

	})
	if err != nil {
		return err
	}
	defer sub.Unsubscribe()

	<-ctx.Done()
	log.Println("ride request service: exiting")
	return nil
}
