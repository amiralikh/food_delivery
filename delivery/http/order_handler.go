package http

import (
	"encoding/json"
	"foodDelivery/domain"
	"foodDelivery/usecase"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
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

func (oh *OrderHandler) GetUserOrders(w http.ResponseWriter, r *http.Request) {
	userID := 7
	orders, err := oh.orderUseCase.GetUserOrders(int64(userID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Serialize the foods object into JSON.
	response, err := json.Marshal(orders)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(response)
}

func (oh *OrderHandler) GetOrderWithItems(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderIDStr := vars["id"]
	orderID, _ := strconv.ParseInt(orderIDStr, 10, 64)
	order, err := oh.orderUseCase.GetOrderWithItems(orderID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Serialize the foods object into JSON.
	response, err := json.Marshal(order)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(response)
}
