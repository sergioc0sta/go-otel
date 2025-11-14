package main

import (
	"log"
	"net/http"

	"github.com/sergioc0sta/go-otel/config"
	"github.com/sergioc0sta/go-otel/internal/infra/handlers"
)

func init(){
	err := config.LoadConfig("./.env")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
		panic(err)
	}
}
func main(){
		mux := http.NewServeMux()

		mux.HandleFunc("/temperature", handlers.TemperatureHandler)

		log.Println("Service B is running on 8080 port...")
		http.ListenAndServe(":8080", mux)
}

