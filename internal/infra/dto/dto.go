package dto

type CepInput struct {
	Cep  string  `json:"cep"`
}

type LocationResponse struct {
	Location string `json:"localidade"`
}

type TemperatureResponse struct {
	City   string  `json:"city"`
	TempC float64 `json:"temp_C"`
	TempF float64 `json:"temp_F"`
	TempK float64 `json:"temp_K"`
}

type TemperatureAPIResponse struct {
	Results struct {
		Temp int `json:"temp"`
	} `json:"results"`
}
