package usecase

import (
	"foodDelivery/domain"
	"foodDelivery/repository"
)

// UserUseCase represents the user use case interface.
type UserUseCase interface {
	GetUserByID(userID int64) (*domain.User, error)
	GetUserByEmail(userEmail string) (*domain.User, error)
	RegisterUser(user *domain.User) error
	CreateUser(user *domain.User) error
	UpdateUser(user *domain.User) error
	DeleteUser(userID int64) error
}

// userUseCase represents the user use case implementation.
type userUseCase struct {
	userRepository repository.UserRepository
}

// NewUserUseCase creates a new instance of UserUseCase.
func NewUserUseCase(userRepository repository.UserRepository) UserUseCase {
	return &userUseCase{
		userRepository: userRepository,
	}
}

// GetUserByID retrieves a user by ID.
func (uc *userUseCase) GetUserByID(userID int64) (*domain.User, error) {
	user, err := uc.userRepository.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (uc *userUseCase) GetUserByEmail(userEmail string) (*domain.User, error) {
	user, err := uc.userRepository.GetUserByEmail(userEmail)
	if err != nil {
		return nil, repository.ErrUserNotFound
	}

	return user, nil
}

// CreateUser creates a new user.
func (uc *userUseCase) CreateUser(user *domain.User) error {
	err := uc.userRepository.CreateUser(user)
	if err != nil {
		return err
	}

	return nil
}

func (uc *userUseCase) RegisterUser(user *domain.User) error {
	err := uc.userRepository.RegisterUser(user)
	if err != nil {
		return err
	}

	return nil
}

// UpdateUser updates an existing user.
func (uc *userUseCase) UpdateUser(user *domain.User) error {
	err := uc.userRepository.UpdateUser(user)
	if err != nil {
		return err
	}

	return nil
}

// DeleteUser deletes a user.
func (uc *userUseCase) DeleteUser(userID int64) error {
	err := uc.userRepository.DeleteUser(userID)
	if err != nil {
		return err
	}

	return nil
}
