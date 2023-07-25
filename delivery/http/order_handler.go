package http

import (
	"encoding/json"
	"foodDelivery/domain"
	"foodDelivery/usecase"
	"net/http"
)

type OrderHandler struct {
	orderUseCase usecase.OrderUseCase
}

func NewOrderHandler(orderUseCase usecase.OrderUseCase) *OrderHandler {
	return &OrderHandler{
		orderUseCase: orderUseCase,
	}
}

func (oh *OrderHandler) SubmitOrder(w http.ResponseWriter, r *http.Request) {
	var order domain.Order
	err := json.NewDecoder(r.Body).Decode(&order)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = oh.orderUseCase.SubmitOrder(&order)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	response := []byte(`{"message": "Order Placed successfully"}`)
	_, _ = w.Write(response)
}
