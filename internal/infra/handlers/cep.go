package handlers

import (
	"fmt"
	"io"
	"net/http"
)

func CepHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Erro ao ler body: %v", err)
		return
	}
	defer r.Body.Close()

	fmt.Println("Body:", string(body))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}
