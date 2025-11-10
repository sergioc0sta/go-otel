package temperature

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type TemperatureClient struct {
	URLCEP  string
	URLTEMP string
	HTTP    *http.Client
}

type LocationResponse struct {
	Location string `json:"localidade"`
}

type TemperatureAPIResponse struct {
	Results struct {
		Temp int `json:"temp"`
	} `json:"results"`
}

type TemperatureResponse struct {
	Celsius    float64 `json:"temp_C"`
	Fahrenheit float64 `json:"temp_F"`
	Kelvin     float64 `json:"temp_K"`
}

func NewTemperatureClient(urlCep, urlTemp string) *TemperatureClient {
	return &TemperatureClient{
		URLCEP:  urlCep,
		URLTEMP: urlTemp,
		HTTP: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (vc *TemperatureClient) getJSON(ctx context.Context, fullURL string, out interface{}) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fullURL, nil)
	if err != nil {
		return err
	}

	resp, err := vc.HTTP.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	return json.NewDecoder(resp.Body).Decode(out)
}


func (vc *TemperatureClient) FetchTemperature(ctx context.Context, cep string) (*TemperatureResponse, error) {

	var data LocationResponse
	err := vc.getJSON(ctx, vc.URLCEP+cep+"/json/", &data)

	if err != nil {
		return nil, err
	}

	escapedLoc := url.PathEscape(data.Location)
	tempURL := fmt.Sprintf("%s%s", vc.URLTEMP, escapedLoc)

	var temperatureAPI TemperatureAPIResponse
	err = vc.getJSON(ctx, tempURL, &temperatureAPI)

	if err != nil {
		return nil, err
	}


	temps := &TemperatureResponse{
		Celsius:    float64(temperatureAPI.Results.Temp),
		Fahrenheit: float64(temperatureAPI.Results.Temp)*9/5 + 32,
		Kelvin:     float64(temperatureAPI.Results.Temp) + 273.15,
	}

	return temps, nil
}
