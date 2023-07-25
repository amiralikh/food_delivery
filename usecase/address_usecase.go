package usecase

import (
	"errors"
	"foodDelivery/domain"
	"foodDelivery/repository"
)

var (
	ErrAddressValidation = errors.New("all fields are required")
)

type AddressUseCase interface {
	CreateAddress(address *domain.Address) error
	UpdateAddress(addressID int64, address *domain.Address) error
	GetAddressByID(addressID int64) (*domain.Address, error)
	DeleteAddress(userID int64, addressID int64) error
	GetUsersAddresses(userID int64) ([]*domain.Address, error)
}

type addressUseCase struct {
	addressRepo repository.AddressRepository
}

func NewAddressUseCase(addressRepo repository.AddressRepository) AddressUseCase {
	return &addressUseCase{
		addressRepo: addressRepo,
	}
}

func (au *addressUseCase) CreateAddress(address *domain.Address) error {
	if address.UserID <= 0 {
		return ErrAddressValidation
	}
	if address.Name == "" {
		return ErrAddressValidation
	}
	if address.Zip == "" {
		return ErrAddressValidation
	}
	if address.Phone == "" {
		return ErrAddressValidation
	}
	if address.Address == "" {
		return ErrAddressValidation
	}

	err := au.addressRepo.CreateAddress(address)
	if err != nil {
		return err
	}

	return nil
}

func (au *addressUseCase) UpdateAddress(addressID int64, address *domain.Address) error {
	if address.UserID <= 0 {
		return ErrAddressValidation
	}
	if address.Name == "" {
		return ErrAddressValidation
	}
	if address.Zip == "" {
		return ErrAddressValidation
	}
	if address.Phone == "" {
		return ErrAddressValidation
	}
	if address.Address == "" {
		return ErrAddressValidation
	}

	err := au.addressRepo.UpdateAddress(addressID, address)
	if err != nil {
		return err
	}

	return nil
}

func (au *addressUseCase) GetAddressByID(addressID int64) (*domain.Address, error) {
	addresses, err := au.addressRepo.GetAddressByID(addressID)
	if err != nil {
		return nil, err
	}

	return addresses, nil
}

func (au *addressUseCase) DeleteAddress(userID int64, addressID int64) error {
	err := au.addressRepo.DeleteAddress(userID, addressID)
	if err != nil {
		return err
	}

	return nil
}

func (au *addressUseCase) GetUsersAddresses(userID int64) ([]*domain.Address, error) {
	addresses, err := au.addressRepo.GetUsersAddresses(userID)
	if err != nil {
		return nil, err
	}

	if len(addresses) == 0 {
		return nil, repository.ErrAddressNotFound
	}

	return addresses, nil
}
