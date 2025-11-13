package location

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/sergioc0sta/go-otel/config"
)

type LocationClient struct {
	HTTP *http.Client
}

type LocationResponse struct {
	Location string
}

type locationAPIResponse struct {
	Localidade string `json:"localidade"`
	Erro       bool   `json:"erro"`
}

func NewLocationClient() *LocationClient {
	return &LocationClient{
		HTTP: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}	

func (lc *LocationClient) getJSON(ctx context.Context, fullURL string, out interface{}) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fullURL, nil)

	if err != nil {
		return err
	}

	resp, err := lc.HTTP.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	return json.NewDecoder(resp.Body).Decode(out)
}

func (lc *LocationClient) FetchLocation(ctx context.Context, cep string) (*LocationResponse, error) {

	var data locationAPIResponse
	err := lc.getJSON(ctx, config.Cfg.ViaCepAPI+cep+"/json/", &data)

	if err != nil {
		return nil, err
	}

	if data.Erro || data.Localidade == "" {
		return nil, fmt.Errorf("zipcode not found")
	}

	dataResponse := &LocationResponse{
		Location: data.Localidade,
	}

	return dataResponse, nil
}
