package usecase

import (
	"foodDelivery/domain"
	"foodDelivery/repository"
)

type SupplierUseCase interface {
	GetSupplierById(supplierID int64) (*domain.Supplier, error)
	CreateSupplier(supplier *domain.Supplier) error
	UpdateSupplier(supplier *domain.Supplier) error
	DeleteSupplier(supplierID int64) error
	GetAllSuppliers() ([]*domain.Supplier, error)
}

type supplierUseCase struct {
	supplierRepository repository.SupplierRepository
}

func NewSupplierUseCase(supplierRepository repository.SupplierRepository) SupplierUseCase {
	return &supplierUseCase{
		supplierRepository: supplierRepository,
	}
}

func (su *supplierUseCase) GetSupplierById(supplierID int64) (*domain.Supplier, error) {
	supplier, err := su.supplierRepository.GetSupplierByID(supplierID)
	if err != nil {
		return nil, err
	}
	return supplier, nil
}

func (su *supplierUseCase) CreateSupplier(supplier *domain.Supplier) error {
	err := su.supplierRepository.CreateSupplier(supplier)
	if err != nil {
		return err
	}
	return nil
}

func (su *supplierUseCase) UpdateSupplier(supplier *domain.Supplier) error {
	err := su.supplierRepository.UpdateSupplier(supplier)
	if err != nil {
		return err
	}
	return nil
}

func (su *supplierUseCase) DeleteSupplier(supplierID int64) error {
	err := su.supplierRepository.DeleteSupplier(supplierID)
	if err != nil {
		return err
	}
	return nil
}

func (su *supplierUseCase) GetAllSuppliers() ([]*domain.Supplier, error) {
	supplier, err := su.supplierRepository.GetAllSuppliers()
	if err != nil {
		return nil, err
	}
	return supplier, nil
}
