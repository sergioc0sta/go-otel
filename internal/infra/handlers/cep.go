package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/sergioc0sta/go-otel/internal/infra/dto"
	// "github.com/sergioc0sta/go-otel/internal/util"
	"github.com/sergioc0sta/go-otel/internal/validate"
)

func CepHandler(w http.ResponseWriter, r *http.Request) {
	var cepDto dto.CepInput
	var temperatureAPI dto.TemperatureResponse

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	err := json.NewDecoder(r.Body).Decode(&cepDto)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	idValid := validate.CepValidator(cepDto.Cep)

	if !idValid {
		http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
		return
	}

	req, _ := http.NewRequest("GET", "http://localhost:8080/temperature?cep="+cepDto.Cep, nil)
	result, err := client.Do(req)

	if err != nil || result.StatusCode != http.StatusOK{ 
		http.Error(w, "can not find zipcode", http.StatusNotFound)
		return
	}

	defer result.Body.Close()

	json.NewDecoder(result.Body).Decode(&temperatureAPI)
	temperatuteAndLocation, _ := json.Marshal(temperatureAPI)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(temperatuteAndLocation)
}
