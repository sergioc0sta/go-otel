package dto

type CepInput struct {
	Cep  string  `json:"cep"`
}

type LocationResponse struct {
	Location string `json:"localidade"`
}
