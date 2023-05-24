package geocoding

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
	Country   string  `json:"country"`
}

type GeocodingClient struct {
	host     string
	basePath string
	client   http.Client
	token    string
}

const (
	direct  string = "geo/1.0/direct"
	reverse string = "data/2.5/weather"
)
