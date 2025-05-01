package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"nats-us-nai/pkg/tfl"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/nats-io/nats.go"
)

const (
	DefaultRequestePerMinute = 50
)

type TFLRequestGenerator struct {
	RPM       int
	Functions []func(ctx context.Context) error
}

func bikeOccupancyFetchAndPost(api *tfl.API, nc *nats.Conn) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		bikePoints, err := api.BikePoint.GetAllLocations(ctx)
		if err != nil {
			return err
		}

		stub := func(s string) string {
			s = strings.TrimSpace(s)
			s = strings.ToLower(s)
			return strings.ReplaceAll(s, " ", "-")
		}

		for i, bikePoint := range bikePoints {
			log.Println("Processing bike point", i+1, "of", len(bikePoints), ":", bikePoint.CommonName)
			parts := strings.Split(bikePoint.CommonName, ",")

			topic := "bikepoint." + stub(parts[1]) + "." + stub(parts[0])

			var data []byte
			data, err = json.Marshal(bikePoint.MetaData)
			if err != nil {
				return err
			}
			if err = nc.Publish(topic, data); err != nil {
				return fmt.Errorf("error publishing message to NATS: %v", err)
			}
		}
		return nil
	}

}

func (generator *TFLRequestGenerator) Serve(ctx context.Context) {
	delayBetweenRequests := time.Minute / time.Duration(generator.RPM)

	totalFunctions := len(generator.Functions)

	for i := 0; ; i++ {
		select {
		case <-time.After(delayBetweenRequests):
			if err := generator.Functions[i%totalFunctions](ctx); err != nil {
				log.Fatalln(fmt.Errorf("error running function: %v", err))
			}
		case <-ctx.Done():
			return
		}
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	appKey := os.Getenv("TFL_APP_KEY")
	if appKey == "" {
		log.Fatalln("TFL_APP_KEY environment variable is not set")
	}
	client := &http.Client{}
	api := tfl.NewAPI(appKey, client)

	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatalf("Error connecting to NATS: %v", err)
	}
	defer nc.Close()

	f := bikeOccupancyFetchAndPost(api, nc)

	generator := &TFLRequestGenerator{
		RPM:       DefaultRequestePerMinute,
		Functions: []func(ctx context.Context) error{f},
	}

	go generator.Serve(ctx)

	select {
	case <-time.After(5 * time.Minute):
		cancel()
	}
}
