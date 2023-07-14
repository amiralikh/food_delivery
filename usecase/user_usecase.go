package usecase

import (
	_ "errors"

	"foodDelivery/domain"
	"foodDelivery/repository"
)

// UserUseCase represents the use case for user operations.
type UserUseCase struct {
	userRepository repository.UserRepository
}

// NewUserUseCase creates a new instance of UserUseCase.
func NewUserUseCase(userRepository *repository.UserRepository) *UserUseCase {
	return &UserUseCase{
		userRepository: userRepository,
	}
}

// GetUserByID retrieves a user by ID.
func (uc *UserUseCase) GetUserByID(userID int64) (*domain.User, error) {
	user, err := uc.userRepository.FindUserByID(userID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// GetUserByEmail retrieves a user by email.
func (uc *UserUseCase) GetUserByEmail(email string) (*domain.User, error) {
	user, err := uc.userRepository.FindUserByEmail(email)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// GetUserByPhone retrieves a user by phone.
func (uc *UserUseCase) GetUserByPhone(phone string) (*domain.User, error) {
	user, err := uc.userRepository.FindUserByPhone(phone)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// CreateUser creates a new user.
func (uc *UserUseCase) CreateUser(user *domain.User) error {
	// Perform any necessary validation or business logic before creating the user.
	// For example, check if the email or phone is already registered.

	// Hash the user's password before saving it.
	// You can use a package like "golang.org/x/crypto/bcrypt" for password hashing.

	// Call the repository method to create the user.
	err := uc.userRepository.CreateUser(user)
	if err != nil {
		return err
	}

	return nil
}

// UpdateUser updates an existing user.
func (uc *UserUseCase) UpdateUser(user *domain.User) error {
	// Perform any necessary validation or business logic before updating the user.

	// Hash the user's password before updating it, if necessary.

	// Call the repository method to update the user.
	err := uc.userRepository.UpdateUser(user)
	if err != nil {
		return err
	}

	return nil
}

// DeleteUser deletes a user.
func (uc *UserUseCase) DeleteUser(userID int64) error {
	// Perform any necessary validation or business logic before deleting the user.

	// Call the repository method to delete the user.
	err := uc.userRepository.DeleteUser(userID)
	if err != nil {
		return err
	}

	return nil
}
