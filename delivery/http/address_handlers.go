package http

import (
	"encoding/json"
	"errors"
	"foodDelivery/domain"
	"foodDelivery/repository"
	"foodDelivery/usecase"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type AddressHandler struct {
	addressUseCase usecase.AddressUseCase
}

func NewAddressHandler(addressUseCase usecase.AddressUseCase) *AddressHandler {
	return &AddressHandler{
		addressUseCase: addressUseCase,
	}
}

func (ah *AddressHandler) CreateAddress(w http.ResponseWriter, r *http.Request) {
	var address domain.Address
	err := json.NewDecoder(r.Body).Decode(&address)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	userID := 7
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	address.UserID = int64(userID)

	err = ah.addressUseCase.CreateAddress(&address)
	if err != nil {
		http.Error(w, "Failed to create address", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	response := map[string]interface{}{
		"message": "Address created successfully",
		"address": address,
	}
	json.NewEncoder(w).Encode(response)
}

func (ah *AddressHandler) GetUsersAddresses(w http.ResponseWriter, r *http.Request) {
	userID := 7

	//userID, err := getUserIDFromContext(r.Context())
	//if err != nil {
	//	http.Error(w, "Unauthorized", http.StatusUnauthorized)
	//	return
	//}

	addresses, err := ah.addressUseCase.GetUsersAddresses(int64(userID))
	if err != nil {
		http.Error(w, "Failed to fetch addresses", http.StatusInternalServerError)
		return
	}

	if addresses == nil {
		http.Error(w, "No addresses found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(addresses)
}

func (ah *AddressHandler) UpdateAddress(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	addressIDStr := vars["id"]
	addressID, err := strconv.ParseInt(addressIDStr, 10, 64)

	var address domain.Address
	err = json.NewDecoder(r.Body).Decode(&address)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	userID := 7

	//userID, err := getUserIDFromContext(r.Context())
	//if err != nil {
	//	http.Error(w, "Unauthorized", http.StatusUnauthorized)
	//	return
	//}

	address.UserID = int64(userID)

	err = ah.addressUseCase.UpdateAddress(addressID, &address)
	if err != nil {
		http.Error(w, "Failed to update address", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	response := map[string]interface{}{
		"message": "Address updated successfully",
		"address": address,
	}
	json.NewEncoder(w).Encode(response)
}

func (ah *AddressHandler) DeleteAddress(w http.ResponseWriter, r *http.Request) {
	addressIDStr := mux.Vars(r)["id"]
	addressID, err := strconv.ParseInt(addressIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid address ID", http.StatusBadRequest)
		return
	}

	userID := 7

	//userID, err := getUserIDFromContext(r.Context())
	//if err != nil {
	//	http.Error(w, "Unauthorized", http.StatusUnauthorized)
	//	return
	//}

	err = ah.addressUseCase.DeleteAddress(int64(userID), addressID)
	if err != nil {
		http.Error(w, "Failed to delete address", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	response := map[string]string{"message": "Address deleted successfully"}
	json.NewEncoder(w).Encode(response)
}

func (ah *AddressHandler) GetAddressByID(w http.ResponseWriter, r *http.Request) {
	addressIDStr := mux.Vars(r)["id"]
	addressID, err := strconv.ParseInt(addressIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid address ID", http.StatusBadRequest)
		return
	}

	address, err := ah.addressUseCase.GetAddressByID(addressID)
	if err != nil {
		if errors.Is(err, repository.ErrAddressNotFound) {
			http.Error(w, "Address not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to fetch address", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(address)
}
