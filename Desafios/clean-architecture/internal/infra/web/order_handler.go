package web

import (
	"encoding/json"
	"net/http"

	"github.com/Whallas/goexpert/Desafios/clean-architecture/internal/usecase"
)

type OrderHandler struct {
	CreateOrderUseCase usecase.CreateOrderUseCase
	ListOrdersUseCase  usecase.ListOrdersUseCase
}

func NewOrderHandler(create usecase.CreateOrderUseCase, list usecase.ListOrdersUseCase) *OrderHandler {
	return &OrderHandler{
		CreateOrderUseCase: create,
		ListOrdersUseCase:  list,
	}
}

func (h *OrderHandler) Create(w http.ResponseWriter, r *http.Request) {
	var input usecase.OrderInputDTO
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	output, err := h.CreateOrderUseCase.Execute(input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(output)
}

func (h *OrderHandler) List(w http.ResponseWriter, r *http.Request) {
	output, err := h.ListOrdersUseCase.Execute()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(output)
}
