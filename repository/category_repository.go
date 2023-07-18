package repository

import (
	"database/sql"
	"foodDelivery/domain"
)

type CategoryRepository interface {
	GetCategoryByID(categoryID int64) (*domain.Category, error)
	CreateCategory(category *domain.Category) error
	UpdateCategory(category *domain.Category) error
	DeleteCategory(categoryID int64) error
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
	query := "INSERT INTO categories (name, image_url) VALUES ($s1, $s2)"
	_, err := cr.db.Exec(query, category.Name, category.ImageUrl)
	if err != nil {
		return err
	}
	return nil
}

func (cr *categoryRepository) UpdateCategory(category *domain.Category) error {
	query := "UPDATE categories SET name= $s1, image_url = $2"
	_, err := cr.db.Exec(query, category.Name, category.ImageUrl)
	if err != nil {
		return err
	}

	return nil
}

func (cr *categoryRepository) DeleteCategory(categoryID int64) error {
	query := "DELETE FROM categories WHERE id = $1"
	_, err := cr.db.Exec(query, categoryID)
	if err != nil {
		return err
	}
	return nil
}
