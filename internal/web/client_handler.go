package web

import (
	"encoding/json"
	"fmt"
	"github.com/danilobandeira29/ms-wallet/internal/usecase/createclient"
	"log"
	"net/http"
	"os"
)

type WebClientHandler struct {
	CreateClientUseCase createclient.CreateClientUseCase
}

func NewWebClientHandler(c createclient.CreateClientUseCase) *WebClientHandler {
	return &WebClientHandler{
		CreateClientUseCase: c,
	}
}
func (wc *WebClientHandler) CreateClient(w http.ResponseWriter, r *http.Request) {
	var dto createclient.InputDTO
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	output, err := wc.CreateClientUseCase.Execute(dto)
	if err != nil {
		log.SetOutput(os.Stdout)
		e := fmt.Sprintf("error when execute service %v", err)
		log.Println(e)
		http.Error(w, e, http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(output)
	fmt.Println(err)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}
