# CEP Weather System — Summary

## Overview
Two Go services:
- **Service A (Input):** Validates CEP and forwards to Service B.
- **Service B (Orchestrator):** Resolves city by CEP, fetches weather, returns temps in °C/°F/K.

## Payload & Validation
- **POST body (both):**
  ```json
  { "cep": "29902555" }
  ```
- `cep` **must be a string of exactly 8 digits**.

## API
### Service A
- `POST /cep`
  - Invalid CEP → **422**, `{"message":"invalid zipcode"}`
  - Valid CEP → forwards to Service B and **relays** its response/status

### Service B
- `POST /weather-by-cep`
  - **200**
    ```json
    { "city":"São Paulo","temp_C":28.5,"temp_F":83.3,"temp_K":301.6 }
    ```
  - **422** `"invalid zipcode"`
  - **404** `"can not find zipcode"`

## Conversions
- `°F = °C * 9/5 + 32`  
- `K = °C + 273.15`  
*(Round to sensible precision, e.g., 1 decimal.)*

## Tracing (OTel + Zipkin)
- **Distributed tracing** across A → B with W3C propagation (`traceparent`).
- Spans in B: `cep.lookup` and `weather.lookup`.
- Export to Zipkin (e.g., `http://localhost:9411/api/v2/spans`).

### Como rodar localmente
- Suba o Zipkin localmente (`docker run -d -p 9411:9411 openzipkin/zipkin` ou use o serviço definido no `docker-compose`).
- Configure as variáveis de ambiente esperadas (veja `.env`). Cada serviço precisa do próprio `OTEL_SERVICE_NAME`:
  ```bash
  # Terminal 1
  OTEL_SERVICE_NAME=service-b go run ./cmd/service-b

  # Terminal 2
  OTEL_SERVICE_NAME=service-a go run ./cmd/service-a
  ```
- `SERVICE_B_URL` aponta para o endereço usado pela API A para chamar a API B (padrão `http://localhost:8080`).
- `OTEL_EXPORTER_ZIPKIN_ENDPOINT` deve apontar para `http://<zipkin-host>:9411/api/v2/spans`.
- As requisições HTTP de entrada/saída já são instrumentadas via `otelhttp` e o serviço B cria spans dedicados (`cep.lookup` e `weather.lookup`) para as chamadas externas.

## Suggested ENV
- **A:** `PORT`, `SERVICE_B_URL`, `OTEL_EXPORTER_ZIPKIN_ENDPOINT`, `OTEL_SERVICE_NAME`
- **B:** `PORT`, `CEP_PROVIDER_URL`, `WEATHER_PROVIDER_URL`, `OTEL_EXPORTER_ZIPKIN_ENDPOINT`, `OTEL_SERVICE_NAME`

## Local (optional)
- Dockerfiles for both services + `docker-compose` with `service-a`, `service-b`, and `zipkin`.  
