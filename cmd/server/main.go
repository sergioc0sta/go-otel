package main

import (
	"log"
	"net/http"

	"github.com/sergioc0sta/go-otel/config"
	"github.com/sergioc0sta/go-otel/internal/infra/handlers"
	"github.com/sergioc0sta/go-otel/internal/infra/temperature"
)

func init() {
	err := config.LoadConfig("./.env")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
		panic(err)
	}
}

func main() {
	mux := http.NewServeMux()
	temperatureClient := temperature.NewTemperatureClient(config.Cfg.ViaCepAPI, config.Cfg.WeatherAPI)

	mux.HandleFunc("/health", handlers.HealthHandler)
	mux.HandleFunc("/temp", handlers.WeatherHandler(temperatureClient))

	log.Println("Server is running on 8080 port...")
	http.ListenAndServe(":8080", mux)
}
