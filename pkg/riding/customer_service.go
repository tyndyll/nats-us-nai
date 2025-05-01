package riding

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"log"
	"time"
)

type CustomerService struct {
	ID string

	nc *nats.Conn
}

func (cs *CustomerService) Serve(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	rr := &RideRequest{
		RequestID:   uuid.New().String(),
		Destination: GetRandomLocation(),
		CustomerID:  cs.ID,
		TimeStamp:   time.Now().Format(time.RFC3339),
	}

	data, err := json.Marshal(rr)
	if err != nil {
		log.Fatalln("could not marshal ride request:", err)
	}

	// Subscribe to the customer topic
	rideCh := make(chan *nats.Msg, 10)
	// We're only looking for the ride accepted message, but we could use a wildcard here
	topic := CustomerRideAcceptedTopic(cs.ID)
	rideSub, err := cs.nc.ChanSubscribe(topic, rideCh)
	if err != nil {
		log.Fatalln("could not subscribe to ride accepted topic:", err)
	}
	defer rideSub.Unsubscribe()

	log.Println("customer-service: publishing ride request:", string(data))
	if err = cs.nc.Publish(RideRequestTopic, data); err != nil {
		log.Fatalln("could not publish ride request:", err)
	}

	for {
		select {
		case msg := <-rideCh:
			log.Println("customer-service: received ride update:", string(msg.Data))
			cancel()
		case <-ctx.Done():
			log.Println("context done, exiting...")
			return nil
		case <-time.After(30 * time.Second):
			log.Println("customer-service: no response from driver, timer ticking")
			return fmt.Errorf("timed out waiting for ride update")
		}
	}
}

func (cs *CustomerService) HandleRideAccepted(msg *nats.Msg) {
	log.Println("received ride accepted:", string(msg.Data))
}

func NewCustomerService(nc *nats.Conn) *CustomerService {
	return &CustomerService{
		ID: uuid.New().String(),
		nc: nc,
	}
}
