package gcClient

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"

	"main.go/clients/events"
	"main.go/lib/e"
)

func New(host string, token string) *GeocodingClient {
	return &GeocodingClient{
		host:   host,
		client: http.Client{},
		token:  token,
	}
}

func (gc *GeocodingClient) FetchCity(city string) ([]events.Weatherdata, error) {
	cities, err := gc.queryCityName(city)

	if err != nil {
		return nil, e.Wrap("can't fetch city: %w", err)
	}

	return convertToWeatherData(cities), nil
}

func (gc *GeocodingClient) FetchCityWithCoord(lat, lon float32) ([]events.Weatherdata, error) {
	cities, err := gc.queryCoordinates(lat, lon)

	if err != nil {
		return nil, e.Wrap("can't query with coordinates", err)
	}

	return convertToWeatherData(cities), nil
}

func (gc *GeocodingClient) queryCityName(city string) ([]CityStats, error) {
	query := url.Values{}
	query.Add("appid", gc.token)
	query.Add("q", city)

	return gc.processQuery(query, direct)
}

func (gc *GeocodingClient) queryCoordinates(lat, lon float32) ([]CityStats, error) {
	query := url.Values{}
	query.Add("appid", gc.token)
	query.Add("lat", strconv.FormatFloat(float64(lat), 'f', 2, 32))
	query.Add("lon", strconv.FormatFloat(float64(lon), 'f', 2, 32))

	return gc.processQuery(query, reverse)
}

func (gc *GeocodingClient) processQuery(query url.Values, requestType string) ([]CityStats, error) {
	response, err := gc.doRequest(query, requestType)
	if err != nil {
		return nil, e.Wrap("can't do request with query", err)
	}

	var gcResponse GeocodingResponse

	if err := json.Unmarshal(response, &gcResponse.Result); err != nil {
		return nil, e.Wrap("can't unmarshal geocoding response", err)
	}

	return gcResponse.Result, nil
}

func convertToWeatherData(cities []CityStats) []events.Weatherdata {
	var citiesWeather []events.Weatherdata

	for _, city := range cities {
		citiesWeather = append(citiesWeather, weather(city))
	}

	return citiesWeather
}

func weather(city CityStats) events.Weatherdata {
	return events.Weatherdata{
		CityName:  getCityname(city),
		Latitude:  getLatitude(city),
		Longitude: getLongitude(city),
	}
}

func getCityname(city CityStats) string {
	return city.Name
}

func getLongitude(city CityStats) float32 {
	return city.Longitude
}

func getLatitude(city CityStats) float32 {
	return city.Latitude
}

func (gc *GeocodingClient) doRequest(query url.Values, requestType string) ([]byte, error) {
	url := url.URL{
		Scheme:   "https",
		Host:     gc.host,
		Path:     path.Join(gc.basePath, requestType),
		RawQuery: query.Encode(),
	}
	fmt.Println(url.String())

	req, err := http.NewRequest(http.MethodGet, url.String(), nil)
	if err != nil {
		return nil, e.Wrap("can't create http request, geocoding", err)
	}

	resp, err := gc.client.Do(req)
	if err != nil {
		return nil, e.Wrap("can't send http request, geocoding", err)
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, e.Wrap("can't read geocoding response", err)
	}

	return body, nil
}
