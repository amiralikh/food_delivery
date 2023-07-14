package repository

import (
	"database/sql"
	"errors"

	"foodDelivery/domain"
)

// UserRepository represents the repository for user data.
type UserRepository struct {
	db *sql.DB
}

// NewUserRepository creates a new instance of UserRepository.
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

// FindUserByID retrieves a user by ID.
func (ur *UserRepository) FindUserByID(userID int64) (*domain.User, error) {
	user := &domain.User{}

	query := "SELECT id, name, last_name, phone, email, password, status FROM users WHERE id = $1"
	err := ur.db.QueryRow(query, userID).Scan(
		&user.ID, &user.Name, &user.LastName, &user.Phone, &user.Email, &user.Password, &user.Status,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return user, nil
}

// FindUserByEmail retrieves a user by email.
func (ur *UserRepository) FindUserByEmail(email string) (*domain.User, error) {
	user := &domain.User{}

	query := "SELECT id, name, last_name, phone, email, password, status FROM users WHERE email = $1"
	err := ur.db.QueryRow(query, email).Scan(
		&user.ID, &user.Name, &user.LastName, &user.Phone, &user.Email, &user.Password, &user.Status,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return user, nil
}

// FindUserByPhone retrieves a user by phone.
func (ur *UserRepository) FindUserByPhone(phone string) (*domain.User, error) {
	user := &domain.User{}

	query := "SELECT id, name, last_name, phone, email, password, status FROM users WHERE phone = $1"
	err := ur.db.QueryRow(query, phone).Scan(
		&user.ID, &user.Name, &user.LastName, &user.Phone, &user.Email, &user.Password, &user.Status,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return user, nil
}

// CreateUser creates a new user.
func (ur *UserRepository) CreateUser(user *domain.User) error {
	query := "INSERT INTO users(name, last_name, phone, email, password, status) VALUES($1, $2, $3, $4, $5, $6) RETURNING id"
	err := ur.db.QueryRow(query, user.Name, user.LastName, user.Phone, user.Email, user.Password, user.Status).Scan(&user.ID)
	if err != nil {
		return err
	}

	return nil
}

// UpdateUser updates an existing user.
func (ur *UserRepository) UpdateUser(user *domain.User) error {
	query := "UPDATE users SET name = $1, last_name = $2, phone = $3, email = $4, password = $5, status = $6 WHERE id = $7"
	_, err := ur.db.Exec(query, user.Name, user.LastName, user.Phone, user.Email, user.Password, user.Status, user.ID)
	if err != nil {
		return err
	}

	return nil
}

// DeleteUser deletes a user.
func (ur *UserRepository) DeleteUser(userID int64) error {
	query := "DELETE FROM users WHERE id = $1"
	_, err := ur.db.Exec(query, userID)
	if err != nil {
		return err
	}

	return nil
}
