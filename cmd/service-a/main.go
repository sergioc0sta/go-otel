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

		mux.HandleFunc("/cep", handlers.CepHandler)

		log.Println("Service A is running on 8081 port...")
		http.ListenAndServe(":8081", mux)
}

