package web

import (
	"encoding/json"
	"github.com/danilobandeira29/ms-wallet/internal/usecase/createaccount"
	"net/http"
)

type WebAccountHandler struct {
	CreateAccountUseCase createaccount.UseCase
}

func NewWebAccountHandler(c createaccount.UseCase) *WebAccountHandler {
	return &WebAccountHandler{CreateAccountUseCase: c}
}

func (wa *WebAccountHandler) CreateClient(w http.ResponseWriter, r *http.Request) {
	var input createaccount.InputDTO
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	output, err := wa.CreateAccountUseCase.Execute(input)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = json.NewEncoder(w).Encode(output)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}
