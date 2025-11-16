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

### Run Docker Compose
Run  `docker-compose up --build` :
- `service-a` (porta `8081`)
- `service-b` (porta `8080`)
- `zipkin` (porta `9411`)

Instead of if you have the `make` installed
- ```sh
    make run-compose-up
    ```
