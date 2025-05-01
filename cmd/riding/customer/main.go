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
	svc := riding.NewCustomerService(nc)
	if err = svc.Serve(ctx); err != nil {
		log.Println("error running customer service:", err)
	}
}
