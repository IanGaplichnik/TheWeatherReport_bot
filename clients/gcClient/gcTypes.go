package gcClient

import (
	"net/http"
)

type GeocodingResponse struct {
	Result []CityStats
}

type CityStats struct {
	Name      string  `json:"name"`
	Latitude  float32 `json:"lat"`
	Longitude float32 `json:"lon"`
}

type GeocodingClient struct {
	host     string
	basePath string
	client   http.Client
	token    string
}

const (
	direct  string = "geo/1.0/direct"
	reverse string = "geo/1.0/reverse"
)
