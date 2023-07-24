package usecase

import (
	"errors"
	"foodDelivery/domain"
	"foodDelivery/repository"
)

var (
	ErrFoodNameRequired = errors.New("food name is required")
	ErrCategoryRequired = errors.New("food category is required")
	ErrCategoryNotFound = errors.New("food category not found")
	ErrSupplierRequired = errors.New("food supplier is required")
	ErrSupplierNotFound = errors.New("food supplier not found")
	ErrFoodNotFound     = errors.New("food not found")
	ErrImageNotFound    = errors.New("image not found")
)

type FoodUseCase interface {
	CreateFood(food *domain.Food) error
	GetFoodByID(foodID int64) (*domain.Food, error)
	UpdateFood(food *domain.Food) error
	SyncGallery(foodID int64, images []*domain.Image) error
	DeleteFood(foodID int64) error
}

type foodUseCase struct {
	foodRepo     repository.FoodRepository
	categoryRepo repository.CategoryRepository
	supplierRepo repository.SupplierRepository
	galleryRepo  repository.GalleryRepository
}

func NewFoodUseCase(foodRepo repository.FoodRepository, categoryRepo repository.CategoryRepository,
	supplierRepo repository.SupplierRepository, galleryRepo repository.GalleryRepository) *foodUseCase {
	return &foodUseCase{
		foodRepo:     foodRepo,
		categoryRepo: categoryRepo,
		supplierRepo: supplierRepo,
		galleryRepo:  galleryRepo,
	}
}

func (fu *foodUseCase) GetFoodByID(foodID int64) (*domain.Food, error) {
	food, err := fu.foodRepo.GetFoodByID(foodID)
	if err != nil {
		return nil, ErrFoodNotFound
	}

	return food, nil
}

func (fu *foodUseCase) CreateFood(food *domain.Food) error {
	if food.Name == "" {
		return ErrFoodNameRequired
	}

	if food.CategoryID == 0 {
		return ErrCategoryRequired
	}

	if food.SupplierID == 0 {
		return ErrSupplierRequired
	}

	_, err := fu.categoryRepo.GetCategoryByID(food.CategoryID)
	if err != nil {
		return ErrCategoryNotFound
	}

	_, err = fu.supplierRepo.GetSupplierByID(food.SupplierID)
	if err != nil {
		return ErrSupplierNotFound
	}

	err = fu.foodRepo.CreateFood(food)
	if err != nil {
		return err
	}

	return nil
}

func (fu *foodUseCase) UpdateFood(food *domain.Food) error {
	_, err := fu.supplierRepo.GetSupplierByID(food.SupplierID)
	if err != nil {
		return ErrSupplierNotFound
	}

	_, err = fu.categoryRepo.GetCategoryByID(food.CategoryID)
	if err != nil {
		return ErrCategoryNotFound
	}

	err = fu.foodRepo.UpdateFood(food)
	if err != nil {
		return err
	}

	return nil
}

func (fu *foodUseCase) SyncGallery(foodID int64, images []*domain.Image) error {
	_, err := fu.foodRepo.GetFoodByID(foodID)
	if err != nil {
		if errors.Is(err, repository.ErrFoodNotFound) {
			return ErrFoodNotFound
		}
		return err
	}

	existingGallery, err := fu.galleryRepo.GetImagesByFoodID(foodID)
	if err != nil && !errors.Is(err, repository.ErrGalleryNotFound) {
		return err
	}

	// Create a map to store existing images in the gallery.
	existingGalleryMap := make(map[int64]*domain.Image)
	for _, img := range existingGallery {
		existingGalleryMap[img.ID] = img
	}

	for _, newImg := range images {
		if newImg.ID == 0 {
			// If the image ID is 0, it means it's a new image and needs to be inserted.
			newImg.FoodID = foodID // Set the correct food ID for the new image.
			err := fu.galleryRepo.InsertImage(newImg)
			if err != nil {
				return err
			}
		} else {
			// If the image ID is not 0, it means it's an existing image and needs to be updated.
			existingImg, found := existingGalleryMap[newImg.ID]
			if !found {
				return ErrImageNotFound
			}

			existingImg.ImageURL = newImg.ImageURL
			err := fu.galleryRepo.UpdateImage(existingImg)
			if err != nil {
				return err
			}

			// Remove the updated image from the map to handle deletion later.
			delete(existingGalleryMap, newImg.ID)
		}
	}

	// Delete any remaining images in the existing gallery map, which means they were not part of the update.
	for _, imgToDelete := range existingGalleryMap {
		err := fu.galleryRepo.DeleteImage(imgToDelete.ID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (fu *foodUseCase) DeleteAllImagesByFoodID(foodID int64) error {
	return fu.galleryRepo.DeleteAllImagesByFoodID(foodID)
}

func (fu *foodUseCase) DeleteFood(foodID int64) error {
	_, err := fu.foodRepo.GetFoodByID(foodID)
	if err != nil {
		if errors.Is(err, repository.ErrFoodNotFound) {
			return ErrFoodNotFound
		}
		return err
	}

	// Delete all images associated with the food.
	hasImages, err := fu.galleryRepo.HasImages(foodID)
	if hasImages {
		err = fu.galleryRepo.DeleteAllImagesByFoodID(foodID)
		if err != nil {
			return err
		}
	}
	if err != nil {
		return err
	}

	// Delete the food after its images have been deleted.
	err = fu.foodRepo.DeleteFood(foodID)
	if err != nil {
		return err
	}

	return nil
}
