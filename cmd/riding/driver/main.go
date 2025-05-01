package main

import (
	"context"
	"log"
	"nats-us-nai/pkg/riding"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	nc, err := riding.NewNATSConnection("")
	if err != nil {
		log.Fatalln(err)
	}

	svc := riding.NewDriverService(nc)
	go svc.Serve(ctx)

	<-ctx.Done()
}
