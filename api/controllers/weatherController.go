package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"server/api/controllers/helpers"
	"server/common/appError"
	"server/common/services"
)

func GetWeather(w http.ResponseWriter, r *http.Request) {
	lat, err := helpers.GetQueryParam(r, "lat")
	if err != nil || lat == "" {
		http.Error(w, "Missing lat", http.StatusBadRequest)
		return
	}

	lon, err := helpers.GetQueryParam(r, "lon")
	if err != nil || lon == "" {
		http.Error(w, "Missing lon", http.StatusBadRequest)
		return
	}

	weather, err := services.GetWeatherByCoordinates(parseFloat(lat), parseFloat(lon))
	if err != nil {
		appError.HandleError(w, err)
		return
	}

	json.NewEncoder(w).Encode(weather)
}

func parseFloat(s string) float64 {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}
	return f
}
