package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/sergioc0sta/go-otel/config"
	"github.com/sergioc0sta/go-otel/internal/infra/dto"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

func TemperatureHandler(w http.ResponseWriter, r *http.Request) {
	var location dto.LocationResponse
	var temperatureAPI dto.TemperatureAPIResponse

	fmt.Println("Chamado B")
	tracer := otel.Tracer("service-b")

	client := &http.Client{
		Timeout:   30 * time.Second,
		Transport: otelhttp.NewTransport(http.DefaultTransport),
	}
	cep := r.URL.Query().Get("cep")
	if cep == "" {
		http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
		return
	}

	viaCepBase := strings.TrimRight(config.Cfg.ViaCepAPI, "/")
	cepURL := fmt.Sprintf("%s/%s/json/", viaCepBase, cep)
	cepCtx, cepSpan := tracer.Start(r.Context(), "cep.lookup")
	cepSpan.SetAttributes(attribute.String("zipcode", cep))

	defer cepSpan.End()

	rp, err := http.NewRequestWithContext(cepCtx, http.MethodGet, cepURL, nil)

	if err != nil {
		cepSpan.RecordError(err)
		http.Error(w, "can not find zipcode", http.StatusNotFound)
		return
	}
	response, err := client.Do(rp)

	if err != nil {
		cepSpan.RecordError(err)
		http.Error(w, "can not find zipcode", http.StatusNotFound)
		return
	}

	defer response.Body.Close()

	err = json.NewDecoder(response.Body).Decode(&location)

	if err != nil || location.Location == "" {
		if err != nil {
			cepSpan.RecordError(err)
		}
		http.Error(w, "can not find zipcode", http.StatusNotFound)
		return
	}

	tempURL := fmt.Sprintf("%s%s", config.Cfg.WeatherAPI,url.PathEscape(location.Location))
	weatherCtx, weatherSpan := tracer.Start(r.Context(), "weather.lookup")
	weatherSpan.SetAttributes(attribute.String("city", location.Location))

	defer weatherSpan.End()

	rpp, err := http.NewRequestWithContext(weatherCtx, http.MethodGet, tempURL, nil)

	if err != nil {
		weatherSpan.RecordError(err)
		http.Error(w, "can not find temperature", http.StatusNotFound)
		return
	}

	responsep, err := client.Do(rpp)

	if err != nil {
		weatherSpan.RecordError(err)
		http.Error(w, "can not find temperature", http.StatusNotFound)
		return
	}

	defer responsep.Body.Close()

	err = json.NewDecoder(responsep.Body).Decode(&temperatureAPI)

	if err != nil {
		weatherSpan.RecordError(err)
		http.Error(w, "can not find temperature", http.StatusNotFound)
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
