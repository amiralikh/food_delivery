package http

import (
	"encoding/json"
	"foodDelivery/domain"
	"foodDelivery/usecase"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type SupplierHandler struct {
	supplierUseCase usecase.SupplierUseCase
}

func NewSupplierHandler(supplierUseCase usecase.SupplierUseCase) *SupplierHandler {
	return &SupplierHandler{
		supplierUseCase: supplierUseCase,
	}
}

func (sh *SupplierHandler) GetSupplierByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	SupplierIDStr := vars["id"]
	SupplierID, _ := strconv.ParseInt(SupplierIDStr, 10, 64)

	supplier, err := sh.supplierUseCase.GetSupplierById(SupplierID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response, _ := json.Marshal(supplier)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(response)
}

func (sh *SupplierHandler) CreateSupplier(w http.ResponseWriter, r *http.Request) {
	var supplier domain.Supplier
	err := json.NewDecoder(r.Body).Decode(&supplier)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	err = sh.supplierUseCase.CreateSupplier(&supplier)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	response := []byte(`{"message": "Supplier created successfully"}`)
	_, _ = w.Write(response)
}

func (sh *SupplierHandler) UpdateSupplier(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	SupplierIDStr := vars["id"]
	SupplierID, _ := strconv.ParseInt(SupplierIDStr, 10, 64)
	var supplier domain.Supplier

	err := json.NewDecoder(r.Body).Decode(&supplier)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	supplier.ID = SupplierID
	err = sh.supplierUseCase.UpdateSupplier(&supplier)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := []byte(`{"message": "Supplier updated successfully"}`)
	_, _ = w.Write(response)
}

func (sh *SupplierHandler) DeleteSupplier(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	SupplierIDStr := vars["id"]
	SupplierID, _ := strconv.ParseInt(SupplierIDStr, 10, 64)
	err := sh.supplierUseCase.DeleteSupplier(SupplierID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := []byte(`{"message": "Supplier deleted successfully"}`)
	_, _ = w.Write(response)
}

func (sh *SupplierHandler) GetAllSuppliers(w http.ResponseWriter, r *http.Request) {
	suppliers, err := sh.supplierUseCase.GetAllSuppliers()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Serialize the categories into JSON.
	response, err := json.Marshal(suppliers)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(response)
}
