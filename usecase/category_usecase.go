package usecase

import (
	"foodDelivery/domain"
	"foodDelivery/repository"
)

type CategoryUseCase interface {
	GetCategoryByID(categoryID int64) (*domain.Category, error)
	CreateCategory(category *domain.Category) error
	UpdateCategory(category *domain.Category) error
	DeleteCategory(category int64) error
}

type categoryUseCase struct {
	categoryRepository repository.CategoryRepository
}

func NewCategoryUseCase(categoryRepository repository.CategoryRepository) CategoryUseCase {
	return &categoryUseCase{
		categoryRepository: categoryRepository,
	}
}

func (cu *categoryUseCase) GetCategoryByID(categoryID int64) (*domain.Category, error) {
	category, err := cu.categoryRepository.GetCategoryByID(categoryID)
	if err != nil {
		return nil, err
	}
	return category, nil
}

func (cu *categoryUseCase) CreateCategory(category *domain.Category) error {
	err := cu.categoryRepository.CreateCategory(category)
	if err != nil {
		return err
	}

	return nil
}

func (cu *categoryUseCase) UpdateCategory(category *domain.Category) error {
	err := cu.categoryRepository.UpdateCategory(category)
	if err != nil {
		return err
	}
	return nil
}

func (cu *categoryUseCase) DeleteCategory(categoryID int64) error {
	err := cu.categoryRepository.DeleteCategory(categoryID)
	if err != nil {
		return err
	}
	return nil
}
