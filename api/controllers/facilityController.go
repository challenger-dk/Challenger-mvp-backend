package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"server/api/controllers/helpers"
	"server/common/appError"
	"server/common/dto"
	"server/common/services"
)

func GetFacilities(w http.ResponseWriter, r *http.Request) {
	searchQuery := r.URL.Query().Get("search_query")
	if searchQuery == "" {
		searchQuery = r.URL.Query().Get("q")
	}
	city := r.URL.Query().Get("city")

	limit := 50
	if l := r.URL.Query().Get("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			limit = parsed
		}
	}
	offset := 0
	if o := r.URL.Query().Get("offset"); o != "" {
		if parsed, err := strconv.Atoi(o); err == nil && parsed >= 0 {
			offset = parsed
		}
	}

	var minLat, maxLat, minLon, maxLon float64
	var hasMinLat, hasMaxLat, hasMinLon, hasMaxLon bool
	if v := r.URL.Query().Get("min_lat"); v != "" {
		if parsed, err := strconv.ParseFloat(v, 64); err == nil {
			minLat, hasMinLat = parsed, true
		}
	}
	if v := r.URL.Query().Get("max_lat"); v != "" {
		if parsed, err := strconv.ParseFloat(v, 64); err == nil {
			maxLat, hasMaxLat = parsed, true
		}
	}
	if v := r.URL.Query().Get("min_lon"); v != "" {
		if parsed, err := strconv.ParseFloat(v, 64); err == nil {
			minLon, hasMinLon = parsed, true
		}
	}
	if v := r.URL.Query().Get("max_lon"); v != "" {
		if parsed, err := strconv.ParseFloat(v, 64); err == nil {
			maxLon, hasMaxLon = parsed, true
		}
	}

	facilities, err := services.ListFacilities(searchQuery, city, limit, offset, minLat, maxLat, minLon, maxLon, hasMinLat, hasMaxLat, hasMinLon, hasMaxLon)
	if err != nil {
		appError.HandleError(w, err)
		return
	}

	response := make([]dto.FacilityResponseDto, len(facilities))
	for i, f := range facilities {
		response[i] = dto.ToFacilityResponseDto(f)
	}

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		appError.HandleError(w, err)
		return
	}
}

func GetFacility(w http.ResponseWriter, r *http.Request) {
	id, err := helpers.GetParamId(r)
	if err != nil {
		appError.HandleError(w, err)
		return
	}

	facility, err := services.GetFacilityByID(id)
	if err != nil {
		appError.HandleError(w, err)
		return
	}

	response := dto.ToFacilityResponseDto(facility)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		appError.HandleError(w, err)
		return
	}
}
