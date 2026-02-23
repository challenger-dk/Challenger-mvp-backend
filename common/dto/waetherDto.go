package dto

type WeatherResponse struct {
	Temperature float64 `json:"temperature"`
	Condition   string  `json:"condition"`
	IconURL     string  `json:"icon_url"`
}

type WeatherAPIResponse struct {
	Current struct {
		TempC     float64 `json:"temp_c"`
		Condition struct {
			Text string `json:"text"`
			Icon string `json:"icon"`
		} `json:"condition"`
	} `json:"current"`
}
