package tfl

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	bikePointPath = "BikePoint"
)

type BikePointService service

func (s *BikePointService) GetAllLocations(ctx context.Context) ([]*BikePoint, error) {
	u := fmt.Sprintf("%s/", bikePointPath)

	req, err := s.client.NewRequest(http.MethodGet, u)
	if err != nil {
		return nil, err
	}

	bikePoints := new([]*BikePoint)

	_, err = s.client.Do(ctx, req, bikePoints)
	if err != nil {
		return nil, err
	}

	return *bikePoints, nil
}

type BikePoint struct {
	ID           string   `json:"id"`
	URL          string   `json:"url"`
	CommonName   string   `json:"commonName"`
	PlaceType    string   `json:"placeType"`
	MetaData     Metadata `json:"additionalProperties"`
	Children     []string `json:"children"`
	ChildrenURLs []string `json:"childrenUrls"`
	Latitude     float64  `json:"lat"`
	Longitude    float64  `json:"lon"`
}

type Metadata map[string]string

func (m *Metadata) UnmarshalJSON(data []byte) error {
	var props []additionalProperty
	if err := json.Unmarshal(data, &props); err != nil {
		return err
	}

	values := map[string]string{}

	for _, prop := range props {
		values[prop.Key] = prop.Value
	}
	*m = values
	return nil
}

type additionalProperty struct {
	Category        string `json:"category"`
	Key             string `json:"key"`
	SourceSystemKey string `json:"sourceSystemKey"`
	Value           string `json:"value"`
	Modified        string `json:"modified"`
}
