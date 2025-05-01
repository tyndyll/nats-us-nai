package tfl

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"sync"
)

const (
	BaseAPIURL       = "https://api.tfl.gov.uk/"
	DefaultUserAgent = "go-tfl/0.1"

	AppKeyQueryParam = "app_key"
)

type API struct {
	client *http.Client

	AppKey string

	BaseURL   *url.URL
	UserAgent string

	common service
	// Services used for talking to different parts of the TFL API.
	BikePoint *BikePointService
	Occupancy *OccupancyService

	mtx sync.Mutex
}

func (api *API) Do(ctx context.Context, req *http.Request, v interface{}) (*http.Response, error) {
	req = req.WithContext(ctx)

	resp, err := api.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	switch v := v.(type) {
	case nil:
	case io.Writer:
		_, err = io.Copy(v, resp.Body)
	default:
		decErr := json.NewDecoder(resp.Body).Decode(v)
		if decErr == io.EOF {
			decErr = nil // ignore EOF errors caused by empty response body
		}
		if decErr != nil {

			err = decErr
		}
	}
	return resp, err
}

func (api *API) NewRequest(method, path string) (*http.Request, error) {
	u, err := api.BaseURL.Parse(path)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, u.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", api.UserAgent)
	req.Header.Add("Cache-Control", "no-cache")

	q := req.URL.Query()
	q.Add(AppKeyQueryParam, api.AppKey)
	req.URL.RawQuery = q.Encode()

	return req, nil
}

func (api *API) initialize() {
	if api.BaseURL == nil {
		api.BaseURL, _ = url.Parse(BaseAPIURL)
	}

	if api.UserAgent == "" {
		api.UserAgent = DefaultUserAgent
	}

	api.common.client = api

	api.BikePoint = (*BikePointService)(&api.common)
	api.Occupancy = (*OccupancyService)(&api.common)
}

func NewAPI(appKey string, httpClient *http.Client) *API {
	if httpClient == nil {
		httpClient = &http.Client{}
	}

	cp := *httpClient
	api := &API{
		client:    &cp,
		AppKey:    appKey,
		UserAgent: DefaultUserAgent,
	}
	api.initialize()
	return api
}

type service struct {
	client *API
}
