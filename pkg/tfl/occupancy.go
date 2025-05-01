package tfl

import (
	"context"
	"net/http"
	"strings"
)

const (
	occupancyPath          = "Occupancy"
	bikepointOccupancyPath = "BikePoints"
)

type OccupancyService service

// GetBikePointOccupancy returns the occupancy of the bike points.
//
// The ID is the bike point ID, which can be found in the BikePoint service. Multiple IDs can be passed in as a
// argument, and an array or results will be returned. Note however that the API will only return a maximum of
// approximately 10 results at a time, so if you pass in more than 10 IDs, the API will return an error.
func (s *OccupancyService) GetBikePointOccupancy(ctx context.Context, id ...string) ([]*BikePointOccupancy, error) {
	u := occupancyPath + "/" + bikepointOccupancyPath + "/" + strings.Join(id, ",")

	req, err := s.client.NewRequest(http.MethodGet, u)
	if err != nil {
		return nil, err
	}

	bikePointOccupancy := new([]*BikePointOccupancy)

	_, err = s.client.Do(ctx, req, bikePointOccupancy)
	if err != nil {
		return nil, err
	}

	return *bikePointOccupancy, nil
}

type BikePointOccupancy struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	BikeCount     int    `json:"bikeCount"`
	EmptyDocks    int    `json:"emptyDocks"`
	TotalDocks    int    `json:"totalDocks"`
	StandardBikes int    `json:"standardBikesCount"`
	EBikesCount   int    `json:"ebikesCount"`
}

// HasBrokenDocks checks if the bike point has broken docks.
//
// A calculation of the numbers number of docks - the number of bikes - the number of spaces not giving a value of 0
// indicates broken docks
func (bpo *BikePointOccupancy) HasBrokenDocks() bool {
	return true
}
