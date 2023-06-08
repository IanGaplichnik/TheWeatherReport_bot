package gcClient

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"main.go/lib/e"
)

func (gc *GeocodingClient) FetchWeather(lat, lon float32) (string, error) {
	weatherlist, err := gc.queryWeather(lat, lon)
	if err != nil {
		e.Wrap("can't fetch weather", err)
	}

	var message []byte

	for _, weather := range weatherlist {
		date := strings.Split(weather.Date, " ")
		description := capitalize(weather.Weather[0].Description)
		message = fmt.Appendf(message, "%.5s\n%s\nPrecepitation chance: %.0f%%\n\n",
			date[1],
			description,
			weather.PoP*100)
	}
	return string(message), nil
}

func capitalize(str string) string {
	if len(str) == 0 {
		return ""
	}
	firstChar := str[0]
	firstChar -= 'a' - 'A'

	return string(firstChar) + str[1:]
}

func (gc *GeocodingClient) queryWeather(lat, lon float32) ([]WeatherList, error) {
	query := url.Values{}
	query.Add("appid", gc.token)
	query.Add("cnt", "4")
	query.Add("lat", strconv.FormatFloat(float64(lat), 'f', 2, 32))
	query.Add("lon", strconv.FormatFloat(float64(lon), 'f', 2, 32))

	return gc.processWeatherQuery(query)
}

func (gc *GeocodingClient) processWeatherQuery(query url.Values) ([]WeatherList, error) {
	response, err := gc.doRequest(query, hourlyForecast)
	if err != nil {
		return nil, e.Wrap("can't process weather query", err)
	}

	var weatherdata WeatherForecast

	if err := json.Unmarshal(response, &weatherdata); err != nil {
		return nil, e.Wrap("can't process weather query", err)
	}

	return weatherdata.Result, nil
}
