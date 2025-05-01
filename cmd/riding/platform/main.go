package main

import (
	"context"
	"log"

	"nats-us-nai/pkg/riding"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize NATS connection
	nc, err := riding.NewNATSConnection("")
	if err != nil {
		log.Fatalln(err)
	}
	defer nc.Drain()

	// Create and start the platform service
	svc := riding.NewPlatformService(nc)
	go func() {
		if err = svc.Serve(ctx); err != nil {
			log.Println("error running platform service:", err)
		}
	}()

	// Wait for termination signal
	<-ctx.Done()
}
