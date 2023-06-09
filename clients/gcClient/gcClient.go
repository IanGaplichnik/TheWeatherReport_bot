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

func (gc *GeocodingClient) FetchCitiesByCityName(city string) ([]events.CityData, error) {
	cities, err := gc.queryCitiesByCityName(city)

	if err != nil {
		return nil, e.Wrap("can't fetch city: %w", err)
	}

	return convertToCityData(cities), nil
}

func (gc *GeocodingClient) queryCitiesByCityName(city string) ([]CityStats, error) {
	query := url.Values{}
	query.Add("appid", gc.token)
	query.Add("q", city)
	query.Add("limit", "5")

	cities, err := gc.processQuery(query, coordinatesByName)
	if err != nil {
		return nil, e.Wrap("can't query cities by city name ", err)
	}

	return cities, nil
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

func (gc *GeocodingClient) FetchCityByCoords(lat, lon float32) ([]events.CityData, error) {
	city, err := gc.queryCityByCoords(lat, lon)

	if err != nil {
		return nil, e.Wrap("can't query with coordinates", err)
	}

	return convertToCityData(city), nil
}

func (gc *GeocodingClient) queryCityByCoords(lat, lon float32) ([]CityStats, error) {
	query := url.Values{}
	query.Add("appid", gc.token)
	query.Add("lat", strconv.FormatFloat(float64(lat), 'f', 2, 32))
	query.Add("lon", strconv.FormatFloat(float64(lon), 'f', 2, 32))

	cities, err := gc.processQuery(query, nameByCoordinates)
	if err != nil {
		return nil, e.Wrap("can't query cities by city name ", err)
	}

	return cities, nil
}

func (gc *GeocodingClient) FetchCoordsByCity(city string) (*events.Coordinates, error) {
	cities, err := gc.queryCitiesByCityName(city)

	if err != nil {
		return nil, e.Wrap("can't fetch coords by city", err)
	}

	if len(cities) < 1 {
		return nil, nil
	}

	coords := convertToCoords(cities)

	return coords, nil
}

func convertToCoords(cities []CityStats) *events.Coordinates {
	return &events.Coordinates{
		Latitude:  getLatitude(cities[0]),
		Longitude: getLongitude(cities[0]),
	}
}

func convertToCityData(cities []CityStats) []events.CityData {
	var citiesWeather []events.CityData

	for _, city := range cities {
		citiesWeather = append(citiesWeather, citydata(city))
	}
	fmt.Println("len = " + strconv.Itoa(len(citiesWeather)))

	return citiesWeather
}

func citydata(city CityStats) events.CityData {
	return events.CityData{
		CityName: getCityname(city),
		Country:  getCountry(city),
		State:    getState(city),
	}
}

func getState(city CityStats) string {
	if city.State != nil {
		return *city.State
	}
	return ""
}

func getCountry(city CityStats) string {
	return city.Country
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
