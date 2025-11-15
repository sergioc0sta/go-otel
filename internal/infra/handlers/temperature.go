package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	// "github.com/sergioc0sta/go-otel/config"
	"github.com/sergioc0sta/go-otel/internal/infra/dto"
	// "github.com/sergioc0sta/go-otel/internal/util"
)

	

func TemperatureHandler(w http.ResponseWriter, r *http.Request) {
	var location dto.LocationResponse
	var temperatureAPI dto.TemperatureAPIResponse

	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	cep := r.URL.Query().Get("cep")
	rp, _ := http.NewRequest("GET", "http://viacep.com.br/ws/"+cep+"/json/", nil)
	response, err := client.Do(rp)

	if err != nil {
		http.Error(w, "can find zipcode", http.StatusNotFound)
		return
	}

	defer response.Body.Close()

	err = json.NewDecoder(response.Body).Decode(&location)

	if err != nil {
		http.Error(w, "can find zipcode", http.StatusNotFound)
		return
	}

	escapedLoc := url.PathEscape(location.Location)
	tempURL := fmt.Sprintf("%s%s", "https://api.hgbrasil.com/weather?format=json-cors&city_name=", escapedLoc)
	rpp, err := http.NewRequest("GET", tempURL, nil)
	reponsep, err := client.Do(rpp)

	if err != nil {
		http.Error(w, "can find temperature", http.StatusNotFound)
		return
	}

	defer reponsep.Body.Close()

	err = json.NewDecoder(reponsep.Body).Decode(&temperatureAPI)
	if err != nil {
		http.Error(w, "can find temperature", http.StatusNotFound)
		return
	}

	result := &dto.TemperatureResponse{
		City:  location.Location,
		TempC: float64(temperatureAPI.Results.Temp),
		TempF: float64(temperatureAPI.Results.Temp)*9/5 + 32,
		TempK: float64(temperatureAPI.Results.Temp) + 273.15,
	}

	respBytes, _ := json.Marshal(result)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(respBytes)
}
