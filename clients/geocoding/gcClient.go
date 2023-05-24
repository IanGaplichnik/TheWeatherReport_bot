package geocoding

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"

	"main.go/lib/e"
)

func New(host string, token string) *GeocodingClient {
	return &GeocodingClient{
		host:   host,
		client: http.Client{},
		token:  token,
	}
}

func (gc *GeocodingClient) GeoCoordinates(city string) ([]CityStats, error) {
	query := url.Values{}
	query.Add("appid", gc.token)
	query.Add("q", city)

	response, err := gc.doRequest(query, direct)
	if err != nil {
		return nil, e.Wrap("can't request longitude and latitude", err)
	}

	var gcResponse GeocodingResponse

	if err := json.Unmarshal(response, &gcResponse.Result); err != nil {
		return nil, e.Wrap("can't unmarshal geocoding response", err)
	}

	return gcResponse.Result, nil
}

// func (gc *GeocodingClient) CityWeather(lat, lon float32) (string, error) {
// 	query := url.Values{}
// 	query.Add("appid", gc.token)
// 	query.Add("lat", strconv.FormatFloat(float64(lat), 'f', 2, 32))
// 	query.Add("lon", strconv.FormatFloat(float64(lon), 'f', 2, 32))

// 	response, err := gc.doRequest(query, reverse)
// 	if err != nil {
// 		return "", e.Wrap("can't request weather", err)
// 	}
// }

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

	fmt.Println(string(body[:]) + "\n")
	return body, nil
}
