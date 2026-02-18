package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"server/common/config"
	"server/common/dto"
)

// Weather service integration using WeatherAPI

var baseUrl string = "http://api.weatherapi.com/v1/current.json"

func GetWeatherByCoordinates(lat float64, lon float64) (*dto.WeatherResponse, error) {
	var apiKey string = config.AppConfig.WeatherAPIKey
	if apiKey == "" {
		return nil, fmt.Errorf("Weather API key is not configured")
	}

	url := fmt.Sprintf("%s?key=%s&q=%f,%f", baseUrl, apiKey, lat, lon)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Weather API returned status %d", resp.StatusCode)
	}

	var weatherResponse dto.WeatherAPIResponse
	err = json.NewDecoder(resp.Body).Decode(&weatherResponse)
	if err != nil {
		return nil, err
	}

	weather := dto.WeatherResponse{
		Temperature: weatherResponse.Current.TempC,
		Condition:   weatherResponse.Current.Condition.Text,
		IconURL:     weatherResponse.Current.Condition.Icon,
	}

	return &weather, nil
}
