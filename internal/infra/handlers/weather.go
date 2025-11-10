package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/sergioc0sta/go-otel/internal/infra/temperature"
	"github.com/sergioc0sta/go-otel/internal/validate"
)

type FetchTemp interface {
	FetchTemperature(ctx context.Context, cep string) (*temperature.TemperatureResponse, error)
}

func WeatherHandler(fetchTemp FetchTemp) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		cep := r.URL.Query().Get("cep")
		idValid := validate.CepValidator(cep)

		if !idValid {
			http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
			return
		}

		temperature, err := fetchTemp.FetchTemperature(r.Context(), cep)
		if err != nil || temperature == nil {
			http.Error(w, "can not find zipcode", http.StatusNotFound)
			return
		}

		respBytes, _ := json.Marshal(temperature)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(respBytes)
	}
}
