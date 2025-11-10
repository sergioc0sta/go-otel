package temperature

import (
	"context"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

func TestFetchTemperatureSuccess(t *testing.T) {
	location := "Rio de Janeiro"
	const expectedTemp = 25

	client := NewTemperatureClient("http://mock.example/cep/", "http://mock.example/temp/")
	client.HTTP = &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			switch req.URL.String() {
			case "http://mock.example/cep/01001000/json/":
				return jsonResponse(http.StatusOK, fmt.Sprintf(`{"localidade":%q}`, location)), nil
			case "http://mock.example/temp/" + url.PathEscape(location):
				return jsonResponse(http.StatusOK, fmt.Sprintf(`{"results":{"temp":%d}}`, expectedTemp)), nil
			default:
				t.Fatalf("unexpected request URL: %s", req.URL.String())
				return nil, nil
			}
		}),
	}

	resp, err := client.FetchTemperature(context.Background(), "01001000")
	if err != nil {
		t.Fatalf("FetchTemperature returned error: %v", err)
	}

	if resp == nil {
		t.Fatalf("FetchTemperature returned nil response")
	}

	if resp.Celsius != expectedTemp {
		t.Fatalf("unexpected celsius value: %v", resp.Celsius)
	}

	expectedF := float64(expectedTemp)*9/5 + 32
	if resp.Fahrenheit != expectedF {
		t.Fatalf("unexpected fahrenheit value: %v", resp.Fahrenheit)
	}

	const tolerance = 1e-6
	expectedK := float64(expectedTemp) + 273.15
	if math.Abs(resp.Kelvin-expectedK) > tolerance {
		t.Fatalf("unexpected kelvin value: %v", resp.Kelvin)
	}
}

func TestFetchTemperatureViaCepError(t *testing.T) {
	client := NewTemperatureClient("http://mock.example/cep/", "http://mock.example/temp/")
	client.HTTP = &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			if strings.Contains(req.URL.Path, "/json/") {
				return jsonResponse(http.StatusInternalServerError, ``), nil
			}
			t.Fatalf("weather API should not be called when viaCep fails")
			return nil, nil
		}),
	}

	if _, err := client.FetchTemperature(context.Background(), "01001000"); err == nil {
		t.Fatalf("expected error when viaCep fails")
	}
}

func TestFetchTemperatureWeatherError(t *testing.T) {
	location := "SaoPaulo"
	client := NewTemperatureClient("http://mock.example/cep/", "http://mock.example/temp/")
	client.HTTP = &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			switch req.URL.String() {
			case "http://mock.example/cep/01001000/json/":
				return jsonResponse(http.StatusOK, fmt.Sprintf(`{"localidade":%q}`, location)), nil
			case "http://mock.example/temp/" + url.PathEscape(location):
				return jsonResponse(http.StatusInternalServerError, ``), nil
			default:
				t.Fatalf("unexpected URL: %s", req.URL.String())
				return nil, nil
			}
		}),
	}

	if _, err := client.FetchTemperature(context.Background(), "01001000"); err == nil {
		t.Fatalf("expected error when weather API fails")
	}
}

type roundTripFunc func(req *http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func jsonResponse(status int, body string) *http.Response {
	return &http.Response{
		StatusCode: status,
		Body:       io.NopCloser(strings.NewReader(body)),
		Header:     make(http.Header),
	}
}
