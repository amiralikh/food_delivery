package repository

import (
	"database/sql"
	"errors"
	"foodDelivery/domain"
)

type SupplierRepository interface {
	GetSupplierByID(supplierID int64) (*domain.Supplier, error)
	CreateSupplier(supplier *domain.Supplier) error
	UpdateSupplier(supplier *domain.Supplier) error
	DeleteSupplier(supplierID int64) error
	GetAllSuppliers() ([]*domain.Supplier, error)
}

type supplierRepository struct {
	db *sql.DB
}

func NewSupplierRepository(db *sql.DB) SupplierRepository {
	return &supplierRepository{
		db: db,
	}
}

var (
	ErrSupplierNotFound = errors.New("food not found")
)

func (sr *supplierRepository) GetSupplierByID(supplierID int64) (*domain.Supplier, error) {
	query := "SELECT id,name,address,description,logo_url,opening_hour,closing_hour,user_id,delivery_time FROM suppliers WHERE id = $1"
	row := sr.db.QueryRow(query, supplierID)
	supplier := &domain.Supplier{}
	err := row.Scan(&supplier.ID, &supplier.Name, &supplier.Address, &supplier.Description, &supplier.LogoUrl,
		&supplier.OpeningHour, &supplier.ClosingHour, &supplier.UserID, &supplier.DeliveryTime)
	if err != nil {
		return nil, err
	}
	return supplier, nil
}

func (sr *supplierRepository) CreateSupplier(supplier *domain.Supplier) error {
	query := "INSERT INTO suppliers (name, address, description, logo_url, opening_hour, closing_hour, user_id, delivery_time) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)"
	result, err := sr.db.Exec(query, supplier.Name, supplier.Address, supplier.Description, supplier.LogoUrl, supplier.OpeningHour, supplier.ClosingHour, supplier.UserID, supplier.DeliveryTime)
	rowsAffected, err := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("validation error: supplier not created")
	}
	if err != nil {
		return err
	}

	return nil
}

func (sr *supplierRepository) UpdateSupplier(supplier *domain.Supplier) error {
	query := "UPDATE suppliers SET name = $1, address = $2, description = $3, logo_url = $4, opening_hour = $5, closing_hour = $6, user_id = $7, delivery_time = $8 WHERE id = $9"
	result, err := sr.db.Exec(query, supplier.Name, supplier.Address, supplier.Description, supplier.LogoUrl, supplier.OpeningHour, supplier.ClosingHour, supplier.UserID, supplier.DeliveryTime, supplier.ID)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("validation error: supplier not updated")
	}

	return nil
}

func (sr *supplierRepository) DeleteSupplier(supplierID int64) error {
	query := "DELETE FROM suppliers WHERE id = $1"
	result, err := sr.db.Exec(query, supplierID)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("validation error: supplier not deleted")
	}

	return nil
}

func (sr *supplierRepository) GetAllSuppliers() ([]*domain.Supplier, error) {
	query := "SELECT id,name,address,description,logo_url,opening_hour,closing_hour,user_id,delivery_time FROM suppliers"
	rows, err := sr.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	suppliers := make([]*domain.Supplier, 0)
	for rows.Next() {
		supplier := &domain.Supplier{}
		err := rows.Scan(&supplier.ID, &supplier.Name, &supplier.Address, &supplier.Description, &supplier.LogoUrl,
			&supplier.OpeningHour, &supplier.ClosingHour, &supplier.UserID, &supplier.DeliveryTime)
		if err != nil {
			return nil, err
		}
		suppliers = append(suppliers, supplier)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return suppliers, nil
}
