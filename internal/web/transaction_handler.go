package web

import (
	"encoding/json"
	"github.com/danilobandeira29/ms-wallet/internal/usecase/createtransction"
	"net/http"
)

type WebTransactionHandler struct {
	CreateTransactionUseCase createtransction.UseCase
}

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func NewTransactionHandler(c createtransction.UseCase) *WebTransactionHandler {
	return &WebTransactionHandler{CreateTransactionUseCase: c}
}

func (wt *WebTransactionHandler) CreateClient(w http.ResponseWriter, r *http.Request) {
	var input createtransction.Input
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	ctx := r.Context()
	output, err := wt.CreateTransactionUseCase.Execute(ctx, input)
	if err != nil {
		errStruct := &ErrorResponse{
			Code:    "create_transaction",
			Message: err.Error(),
		}
		errorJson, _ := json.Marshal(errStruct)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(errorJson)
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
