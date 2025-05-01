package riding

import (
	"context"
	"errors"
	"fmt"
	"github.com/nats-io/nats.go"
	"log"
	"sync"
	"time"
)

type DriverInterestService struct {
	nc *nats.Conn

	seen map[string]string
	mtx  sync.Mutex
}

func (dis *DriverInterestService) Serve(ctx context.Context) error {
	log.Println("driver interest service: starting")
	sub, _ := dis.nc.SubscribeSync(DriverInterestedTopic)

	dis.seen = make(map[string]string)

	for {
		m, err := sub.NextMsg(5 * time.Second)
		if err != nil {
			if !errors.Is(nats.ErrTimeout, err) {
				return err
			}
			continue
		}
		log.Printf("driver interest service: Received a driver interest message: %s", string(m.Data))

		msg := &DriverInterested{}
		if err = DecodeMessage(m.Data, msg); err != nil {
			fmt.Printf("Error decoding message: %v\n", err)
			continue
		}

		dis.mtx.Lock()
		driver, found := dis.seen[msg.Request.RequestID]
		if found {
			log.Printf("driver interest service: ride %s already seen by driver %s", msg.Request.RequestID, driver)
			dis.mtx.Unlock()
			continue
		}
		dis.seen[msg.Request.RequestID] = msg.ID
		dis.mtx.Unlock()

		driverTopic := DriverRideAcceptedTopic(msg.ID)
		customerTopic := CustomerRideAcceptedTopic(msg.Request.CustomerID)

		log.Printf("driver interest service: publishing ride %s accepted to driver %s", msg.Request.RequestID, msg.ID)
		if err = dis.nc.Publish(driverTopic, []byte(msg.Request.RequestID)); err != nil {
			log.Printf("Error publishing driver message: %v\n", err)
			continue
		}

		log.Printf("driver interest service: publishing ride %s accepted to customer %s", msg.Request.RequestID, msg.Request.CustomerID)
		if err = dis.nc.Publish(customerTopic, []byte(msg.Request.RequestID)); err != nil {
			log.Printf("Error publishing driver message: %v\n", err)
			continue
		}
	}
}
