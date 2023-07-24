package repository

import (
	"database/sql"
	"errors"
	"foodDelivery/domain"
)

var (
	ErrGalleryNotFound = errors.New("gallery images not found")
)

type GalleryRepository interface {
	CreateImage(image *domain.Image) error
	GetImagesByFoodID(foodID int64) ([]*domain.Image, error)
	UpdateImage(image *domain.Image) error
	DeleteImage(imageID int64) error
	SyncGallery(foodID int64, images []*domain.Image) error
	InsertImage(image *domain.Image) error
}

type galleryRepository struct {
	db *sql.DB
}

func NewGalleryRepository(db *sql.DB) GalleryRepository {
	return &galleryRepository{
		db: db,
	}
}

func (gr *galleryRepository) GetImagesByFoodID(foodID int64) ([]*domain.Image, error) {
	var images []*domain.Image

	rows, err := gr.db.Query("SELECT id, food_id, image_url FROM gallery WHERE food_id = $1", foodID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		image := &domain.Image{}
		err := rows.Scan(&image.ID, &image.FoodID, &image.ImageURL)
		if err != nil {
			return nil, err
		}
		images = append(images, image)
	}

	if len(images) == 0 {
		return nil, ErrGalleryNotFound
	}

	return images, nil
}

func (gr *galleryRepository) CreateImage(image *domain.Image) error {
	var imageID int64

	err := gr.db.QueryRow(`
		INSERT INTO gallery (food_id, image_url)
		VALUES ($1, $2)
		RETURNING id
	`, image.FoodID, image.ImageURL).Scan(&imageID)

	if err != nil {
		return err
	}

	image.ID = imageID
	return nil
}

func (gr *galleryRepository) UpdateImage(image *domain.Image) error {
	_, err := gr.db.Exec(`
		UPDATE gallery
		SET url = $1
		WHERE id = $2
	`, image.ImageURL, image.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("img not found")
		}
		return err
	}
	return nil
}

func (gr *galleryRepository) SyncGallery(foodID int64, images []*domain.Image) error {
	// Step 1: Delete images not in the input slice from the database.
	err := gr.deleteImagesNotInSlice(foodID, images)
	if err != nil {
		return err
	}

	// Step 2: Update or insert images from the input slice into the database.
	err = gr.updateOrInsertImages(foodID, images)
	if err != nil {
		return err
	}

	return nil
}

func (gr *galleryRepository) deleteImagesNotInSlice(foodID int64, images []*domain.Image) error {
	// First, retrieve the IDs of the existing images in the database for the given foodID.
	rows, err := gr.db.Query("SELECT id FROM gallery WHERE food_id = $1", foodID)
	if err != nil {
		return err
	}
	defer rows.Close()

	// Create a map to store the existing image IDs.
	existingImageIDs := make(map[int64]bool)
	for rows.Next() {
		var imageID int64
		if err := rows.Scan(&imageID); err != nil {
			return err
		}
		existingImageIDs[imageID] = true
	}

	// Iterate through the images and delete those that are not present in the existingImageIDs map.
	for _, img := range images {
		if !existingImageIDs[img.ID] {
			_, err = gr.db.Exec("DELETE FROM gallery WHERE id = $1", img.ID)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (gr *galleryRepository) updateOrInsertImages(foodID int64, images []*domain.Image) error {
	println("update or insert", foodID)
	imageIDsMap := make(map[int64]*domain.Image)
	for _, img := range images {
		if img.ID != 0 {
			imageIDsMap[img.ID] = img
		}
	}

	for _, img := range images {
		if img.ID != 0 {
			// Check if the image ID exists in the map.
			if existingImg, ok := imageIDsMap[img.ID]; ok {
				existingImg.ImageURL = img.ImageURL

				_, err := gr.db.Exec(`
					UPDATE gallery
					SET image_url = $1
					WHERE id = $2
				`, existingImg.ImageURL, existingImg.ID)
				if err != nil {
					return err
				}
			}
		} else {
			// If the image ID is 0, it means it's a new image and needs to be inserted.
			img.FoodID = foodID
			err := gr.CreateImage(img)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (gr *galleryRepository) DeleteImage(imageID int64) error {
	res, err := gr.db.Exec("DELETE FROM gallery WHERE id = $1", imageID)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrGalleryNotFound
	}

	return nil
}

func (gr *galleryRepository) InsertImage(image *domain.Image) error {
	// Check if the image has already been assigned an ID.
	if image.ID != 0 {
		return errors.New("insertImage: image ID must be 0 for new images")
	}

	// Insert the new image into the database.
	row := gr.db.QueryRow(`
		INSERT INTO gallery (food_id, image_url)
		VALUES ($1, $2)
		RETURNING id
	`, image.FoodID, image.ImageURL)

	err := row.Scan(&image.ID)
	if err != nil {
		return err
	}

	return nil
}
