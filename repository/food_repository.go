package repository

import (
	"database/sql"
	"errors"
	"foodDelivery/domain"
)

type FoodRepository interface {
	CreateFood(food *domain.Food) error
	GetFoodByID(foodID int64) (*domain.Food, error)
	UpdateFood(food *domain.Food) error
	DeleteFood(foodID int64) error
	GetAllFoodsWithImages() ([]*domain.Food, error)
}

type foodRepository struct {
	db *sql.DB
}

var (
	ErrFoodNotFound = errors.New("food not found")
)

func NewFoodRepository(db *sql.DB) FoodRepository {
	return &foodRepository{
		db: db,
	}
}

func (fr *foodRepository) GetFoodByID(foodID int64) (*domain.Food, error) {
	query := `
		SELECT f.id, f.name, f.supplier_id, s.name AS supplier_name, f.category_id, c.name AS category_name,
			f.image_url, f.description, f.price, f.daily_quantity
		FROM foods f
		INNER JOIN suppliers s ON f.supplier_id = s.id
		INNER JOIN categories c ON f.category_id = c.id
		WHERE f.id = $1
	`
	row := fr.db.QueryRow(query, foodID)
	food := &domain.Food{}
	err := row.Scan(&food.ID, &food.Name, &food.SupplierID, &food.SupplierName, &food.CategoryID, &food.CategoryName,
		&food.ImageUrl, &food.Description, &food.Price, &food.DailyQuantity)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrFoodNotFound
		}
		return nil, err
	}

	galleryRepo := NewGalleryRepository(fr.db)
	food.Gallery, err = galleryRepo.GetImagesByFoodID(food.ID)
	if err != nil {
		return nil, err
	}

	return food, nil
}

func (fr *foodRepository) CreateFood(food *domain.Food) error {
	var foodID int64

	err := fr.db.QueryRow(`
		INSERT INTO foods (name, supplier_id, category_id, image_url, description, price, daily_quantity)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`, food.Name, food.SupplierID, food.CategoryID, food.ImageUrl, food.Description, food.Price, food.DailyQuantity).Scan(&foodID)

	if err != nil {
		return err
	}

	food.ID = foodID

	if len(food.Gallery) > 0 {
		galleryRepo := NewGalleryRepository(fr.db)
		for _, image := range food.Gallery {
			image.FoodID = foodID
			err = galleryRepo.CreateImage(image)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (fr *foodRepository) UpdateFood(food *domain.Food) error {
	_, err := fr.db.Exec(`
		UPDATE foods
		SET name = $1, supplier_id = $2, category_id = $3, image_url = $4, description = $5, price = $6, daily_quantity = $7
		WHERE id = $8
	`, food.Name, food.SupplierID, food.CategoryID, food.ImageUrl, food.Description, food.Price, food.DailyQuantity, food.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrFoodNotFound
		}
		return err
	}

	return nil
}

func (fr *foodRepository) DeleteFood(foodID int64) error {
	_, err := fr.db.Exec("DELETE FROM foods WHERE id = $1", foodID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrFoodNotFound
		}
		return err
	}
	return nil
}

func (fr *foodRepository) GetAllFoodsWithImages() ([]*domain.Food, error) {
	query := `
		SELECT f.id, f.name, f.supplier_id, s.name AS supplier_name, f.category_id, c.name AS category_name,
			f.image_url, f.description, f.price, f.daily_quantity,
			g.id AS image_id, g.image_url AS image_url
		FROM foods f
		INNER JOIN suppliers s ON f.supplier_id = s.id
		INNER JOIN categories c ON f.category_id = c.id
		LEFT JOIN gallery g ON f.id = g.food_id
		ORDER BY f.id
	`

	rows, err := fr.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	foodsMap := make(map[int64]*domain.Food)
	for rows.Next() {
		food := &domain.Food{}
		image := &domain.Image{}
		err := rows.Scan(
			&food.ID, &food.Name, &food.SupplierID, &food.SupplierName, &food.CategoryID, &food.CategoryName,
			&food.ImageUrl, &food.Description, &food.Price, &food.DailyQuantity,
			&image.ID, &image.ImageURL,
		)
		if err != nil {
			return nil, err
		}

		if existingFood, found := foodsMap[food.ID]; found {
			// Append the image to the existing food's gallery.
			existingFood.Gallery = append(existingFood.Gallery, image)
		} else {
			// Create a new food and set its gallery.
			food.Gallery = []*domain.Image{image}
			foodsMap[food.ID] = food
		}
	}

	foods := make([]*domain.Food, 0, len(foodsMap))
	for _, food := range foodsMap {
		foods = append(foods, food)
	}

	return foods, nil
}
