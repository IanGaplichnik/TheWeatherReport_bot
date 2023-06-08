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

type WeatherForecast struct {
	Result []WeatherList `json:"list"`
}

type WeatherList struct {
	Weather []Weather `json:"weather"`
	PoP     float32   `json:"pop"`
	Date    string    `json:"dt_txt"`
}

type Weather struct {
	Description string `json:"description"`
}

const (
	coordinatesByName string = "geo/1.0/direct"
	nameByCoordinates string = "geo/1.0/reverse"
	hourlyForecast    string = "data/2.5/forecast"
)
