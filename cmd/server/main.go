package main

import (
	"log"
	"net/http"

	"github.com/sergioc0sta/go-otel/internal/infra/handlers"
	"github.com/sergioc0sta/go-otel/internal/infra/temperature"
)

const (
	viaCepAPI = "http://viacep.com.br/ws/"
	tempAPI   = "https://api.hgbrasil.com/weather?format=json-cors&city_name="
)

func main() {

	mux := http.NewServeMux()
	temperatureClient := temperature.NewTemperatureClient(viaCepAPI, tempAPI)

	mux.HandleFunc("/health", handlers.HealthHandler)
	mux.HandleFunc("/temp", handlers.WeatherHandler(temperatureClient))

	log.Println("Server is running on 8080 port...")
	http.ListenAndServe(":8080", mux)
}
