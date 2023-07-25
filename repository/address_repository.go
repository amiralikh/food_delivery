package repository

import (
	"database/sql"
	"errors"
	"foodDelivery/domain"
)

type AddressRepository interface {
	CreateAddress(address *domain.Address) error
	GetAddressByID(userID int64) (*domain.Address, error)
	UpdateAddress(addressID int64, address *domain.Address) error
	DeleteAddress(addressID int64, userID int64) error
	GetUsersAddresses(userID int64) ([]*domain.Address, error)
}

type addressRepository struct {
	db *sql.DB
}

var (
	ErrAddressNotFound = errors.New("address not found")
)

func NewAddressRepository(db *sql.DB) AddressRepository {
	return &addressRepository{
		db: db,
	}
}

func (ar *addressRepository) CreateAddress(address *domain.Address) error {
	query := "INSERT INTO addresses (user_id, name, zip, phone, address) VALUES ($1, $2, $3, $4, $5)"
	_, err := ar.db.Exec(query, address.UserID, address.Name, address.Zip, address.Phone, address.Address)
	if err != nil {
		return err
	}
	return nil
}

func (ar *addressRepository) GetAddressByID(addressID int64) (*domain.Address, error) {
	address := &domain.Address{}
	query := "SELECT user_id, name, zip, phone, address FROM addresses WHERE id = $1"
	err := ar.db.QueryRow(query, addressID).Scan(&address.ID, &address.UserID, &address.Name, &address.Zip, &address.Phone, &address.Address)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrAddressNotFound
		}
		return nil, err
	}
	return address, nil
}

func (ar *addressRepository) GetUsersAddresses(userID int64) ([]*domain.Address, error) {
	query := `
		SELECT user_id, name, zip, phone, address
		FROM addresses
		WHERE user_id = $1
	`

	rows, err := ar.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var addresses []*domain.Address
	for rows.Next() {
		var address domain.Address
		err := rows.Scan(&address.ID, &address.UserID, &address.Name, &address.Zip, &address.Phone, &address.Address)
		if err != nil {
			return nil, err
		}
		addresses = append(addresses, &address)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return addresses, nil
}

func (ar *addressRepository) UpdateAddress(addressID int64, address *domain.Address) error {
	query := "UPDATE addresses SET name = $1, zip = $2, phone = $3, address = $4 WHERE user_id = $5 and id = $6"
	_, err := ar.db.Exec(query, address.Name, address.Zip, address.Phone, address.Address, address.UserID, addressID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrAddressNotFound
		}
		return err
	}
	return nil
}

func (ar *addressRepository) DeleteAddress(addressID int64, userID int64) error {
	query := "DELETE FROM addresses WHERE id = $1 AND user_id = $2"
	_, err := ar.db.Exec(query, userID, addressID)
	if err != nil {
		return err
	}
	return nil
}
