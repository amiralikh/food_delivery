package seeder

import (
	"foodDelivery/repository"
	"log"

	"foodDelivery/domain"
	"github.com/bxcodec/faker/v3"
)

type UserSeeder struct {
	userRepository repository.UserRepository
}

func NewUserSeeder(userRepository repository.UserRepository) *UserSeeder {
	return &UserSeeder{
		userRepository: userRepository,
	}
}

func (s *UserSeeder) SeedUsers() {
	var users []domain.User

	for i := 0; i < 50; i++ {
		user := domain.User{
			Name:     faker.FirstName(),
			LastName: faker.LastName(),
			Phone:    faker.Phonenumber(),
			Email:    faker.Email(),
			Password: faker.Password(),
			Status:   "active",
		}
		users = append(users, user)
	}

	for _, user := range users {
		err := s.userRepository.CreateUser(&user)
		if err != nil {
			log.Printf("Failed to create user: %v", err)
		}
	}

	log.Println("User seeding completed")
}
