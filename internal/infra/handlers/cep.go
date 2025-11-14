package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/sergioc0sta/go-otel/config"
	"github.com/sergioc0sta/go-otel/internal/infra/dto"
	"github.com/sergioc0sta/go-otel/internal/util"
	"github.com/sergioc0sta/go-otel/internal/validate"
)

func CepHandler(w http.ResponseWriter, r *http.Request) {
	var cepDto dto.CepInput
	var location dto.LocationResponse

	err := json.NewDecoder(r.Body).Decode(&cepDto)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	idValid := validate.CepValidator(cepDto.Cep)

	if !idValid {
		http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
		return
	}

	err = util.Fetcher(r.Context(), config.Cfg.ViaCepAPI+cepDto.Cep+"/json/", &location)

	if err != nil || location.Location == "" {
		http.Error(w, "timeout", http.StatusRequestTimeout)
		return
	}
	// if I have the location the next step is to call the service B
	respBytes, _ := json.Marshal(location)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(respBytes)
}
