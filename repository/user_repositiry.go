package repository

import (
	"database/sql"
	_ "errors"

	"foodDelivery/domain"
)

// UserRepository represents the user repository interface.
type UserRepository interface {
	GetUserByID(userID int64) (*domain.User, error)
	CreateUser(user *domain.User) error
	UpdateUser(user *domain.User) error
	DeleteUser(userID int64) error
}

// userRepository represents the user repository implementation.
type userRepository struct {
	db *sql.DB
}

// NewUserRepository creates a new instance of UserRepository.
func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

// GetUserByID retrieves a user by ID from the database.
func (ur *userRepository) GetUserByID(userID int64) (*domain.User, error) {
	query := "SELECT id, name, last_name, phone, email, password, status FROM users WHERE id = $1"
	row := ur.db.QueryRow(query, userID)

	user := &domain.User{}
	err := row.Scan(&user.ID, &user.Name, &user.LastName, &user.Phone, &user.Email, &user.Password, &user.Status)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// CreateUser creates a new user in the database.
func (ur *userRepository) CreateUser(user *domain.User) error {
	query := "INSERT INTO users (name, last_name, phone, email, password, status) VALUES ($1, $2, $3, $4, $5, $6)"
	_, err := ur.db.Exec(query, user.Name, user.LastName, user.Phone, user.Email, user.Password, user.Status)
	if err != nil {
		return err
	}

	return nil
}

// UpdateUser updates an existing user in the database.
func (ur *userRepository) UpdateUser(user *domain.User) error {
	query := "UPDATE users SET name = $1, last_name = $2, phone = $3, email = $4, password = $5, status = $6 WHERE id = $7"
	_, err := ur.db.Exec(query, user.Name, user.LastName, user.Phone, user.Email, user.Password, user.Status, user.ID)
	if err != nil {
		return err
	}

	return nil
}

// DeleteUser deletes a user from the database.
func (ur *userRepository) DeleteUser(userID int64) error {
	query := "DELETE FROM users WHERE id = $1"
	_, err := ur.db.Exec(query, userID)
	if err != nil {
		return err
	}

	return nil
}
