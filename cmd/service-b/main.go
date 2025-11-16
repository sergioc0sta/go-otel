package main

import (
	"context"
	"log"
	"net/http"
	"fmt"

	"github.com/sergioc0sta/go-otel/config"
	"github.com/sergioc0sta/go-otel/internal/infra/handlers"
	"github.com/sergioc0sta/go-otel/internal/infra/telemetry"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func init() {
	err := config.LoadConfig("./.env")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
		panic(err)
	}
}
func main() {
	ctx := context.Background()
	shutdown, err := telemetry.SetupProvider(ctx, config.Cfg.ServiceNameB, config.Cfg.OTelExporterEndpoint)

	if err != nil {
		log.Fatalf("failed to setup tracer provider: %v", err)
	}

	defer shutdown(ctx)

	mux := http.NewServeMux()
	mux.Handle("/temperature", otelhttp.NewHandler(http.HandlerFunc(handlers.TemperatureHandler), "temperature.handler"))
	log.Printf("Service B is running on %s port...", config.Cfg.ServiceBPort)
	port :=fmt.Sprintf(":%s", config.Cfg.ServiceBPort)
	http.ListenAndServe(port, mux)
}
