package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/sergioc0sta/go-otel/config"
	"github.com/sergioc0sta/go-otel/internal/infra/dto"
	"github.com/sergioc0sta/go-otel/internal/validate"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func CepHandler(w http.ResponseWriter, r *http.Request) {
	var cepDto dto.CepInput
	var temperatureAPI dto.TemperatureResponse
	client := &http.Client{
		Timeout:   5 * time.Second,
		Transport: otelhttp.NewTransport(http.DefaultTransport),
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

	fullURL := fmt.Sprintf("%s:%s/temperature?cep=%s", config.Cfg.ServiceAPI, config.Cfg.ServiceBPort, cepDto.Cep)
	req, err := http.NewRequestWithContext(r.Context(), http.MethodGet, fullURL, nil)

	if err != nil {
		http.Error(w, "error creating request", http.StatusInternalServerError)
		return
	}

	result, err := client.Do(req)

	if err != nil || result.StatusCode != http.StatusOK {
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
