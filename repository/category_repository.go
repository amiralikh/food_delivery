package repository

import (
	"database/sql"
	"errors"
	"foodDelivery/domain"
)

var (
	ErrCategoryNotFound = errors.New("category not found")
)

type CategoryRepository interface {
	GetCategoryByID(categoryID int64) (*domain.Category, error)
	CreateCategory(category *domain.Category) error
	UpdateCategory(category *domain.Category) error
	DeleteCategory(categoryID int64) error
	GetAllCategories() ([]*domain.Category, error)
	GetCategoriesBySupplierID(supplierID int64) ([]*domain.Category, error)
}

type categoryRepository struct {
	db *sql.DB
}

func NewCategoryRepository(db *sql.DB) CategoryRepository {
	return &categoryRepository{
		db: db,
	}
}

func (cr *categoryRepository) GetCategoryByID(categoryID int64) (*domain.Category, error) {
	query := "SELECT id, name, image_url FROM categories WHERE id=$1"
	row := cr.db.QueryRow(query, categoryID)

	category := &domain.Category{}
	err := row.Scan(&category.ID, &category.Name, &category.ImageUrl)
	if err != nil {
		return nil, err
	}

	return category, nil
}

func (cr *categoryRepository) CreateCategory(category *domain.Category) error {
	if category.Name == "" || category.ImageUrl == "" {
		return errors.New("validation error: category name and image URL must not be empty")
	}

	query := "INSERT INTO categories (name, image_url) VALUES ($1, $2)"
	_, err := cr.db.Exec(query, category.Name, category.ImageUrl)
	if err != nil {
		return err
	}
	return nil
}

func (cr *categoryRepository) UpdateCategory(category *domain.Category) error {
	query := "UPDATE categories SET name = $1, image_url = $2 WHERE id = $3"
	result, err := cr.db.Exec(query, category.Name, category.ImageUrl, category.ID)
	rowsAffected, err := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("validation error: category not found")
	}
	if err != nil {
		return err
	}

	return nil
}

func (cr *categoryRepository) DeleteCategory(categoryID int64) error {
	query := "DELETE FROM categories WHERE id = $1"
	result, err := cr.db.Exec(query, categoryID)
	rowsAffected, err := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("validation error: category not found")
	}
	if err != nil {
		return err
	}
	return nil
}

func (cr *categoryRepository) GetAllCategories() ([]*domain.Category, error) {
	query := "SELECT id, name, image_url FROM categories"
	rows, err := cr.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categories := make([]*domain.Category, 0)
	for rows.Next() {
		category := &domain.Category{}
		err := rows.Scan(&category.ID, &category.Name, &category.ImageUrl)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return categories, nil
}

func (cr *categoryRepository) GetCategoriesBySupplierID(supplierID int64) ([]*domain.Category, error) {
	var categories []*domain.Category
	query := "SELECT id, name, image_url FROM categories WHERE id IN (SELECT DISTINCT category_id FROM foods WHERE supplier_id = $1)"
	rows, err := cr.db.Query(query, supplierID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var category domain.Category
		err := rows.Scan(&category.ID, &category.Name, &category.ImageUrl)
		if err != nil {
			return nil, err
		}
		categories = append(categories, &category)
	}

	return categories, nil
}
