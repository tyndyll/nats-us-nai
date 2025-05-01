package riding

import (
	"context"
	"log"

	"github.com/nats-io/nats.go"
	"github.com/thejerf/suture/v4"
)

type Platform struct {
	nc *nats.Conn

	supervisor *suture.Supervisor
}

func (ps *Platform) Serve(ctx context.Context) error {
	ps.supervisor = suture.NewSimple("Platform Supervisor")

	ps.supervisor.Add(&RideRequestService{
		nc: ps.nc,
	})
	ps.supervisor.Add(&DriverInterestService{
		nc: ps.nc,
	})

	log.Println("platform service: starting")
	return ps.supervisor.Serve(ctx)
}

func NewPlatformService(nc *nats.Conn) *Platform {
	return &Platform{
		nc: nc,
	}
}
