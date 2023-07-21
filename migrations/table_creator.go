package migrations

import (
	"database/sql"
	"log"
)

// CreateUsersTable creates the users table if it doesn't exist.
func CreateUsersTable(db *sql.DB) error {
	// Check if the users table exists.
	tableExists := false
	err := db.QueryRow("SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'users')").Scan(&tableExists)
	if err != nil {
		return err
	}

	// Create the users table if it doesn't exist.
	if !tableExists {
		createUsersTable := `
			CREATE TABLE users (
				id SERIAL PRIMARY KEY,
				name VARCHAR(255) NOT NULL,
				last_name VARCHAR(255) NOT NULL,
				phone VARCHAR(20) NOT NULL,
				email VARCHAR(255) NOT NULL,
				password VARCHAR(255) NOT NULL,
				status VARCHAR(50) NOT NULL
			)
		`

		_, err = db.Exec(createUsersTable)
		if err != nil {
			return err
		}

		log.Println("Users table created successfully")
	} else {
		log.Println("Users table already exists")
	}

	return nil
}

func CreateCategoriesTable(db *sql.DB) error {
	// Check if the users table exists.
	tableExists := false
	err := db.QueryRow("SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'categories')").Scan(&tableExists)
	if err != nil {
		return err
	}

	// Create the users table if it doesn't exist.
	if !tableExists {
		createUsersTable := `
			CREATE TABLE categories (
				id SERIAL PRIMARY KEY,
				name VARCHAR(255) NOT NULL,
				image_url VARCHAR(255) NOT NULL
			)
		`

		_, err = db.Exec(createUsersTable)
		if err != nil {
			return err
		}

		log.Println("category table created successfully")
	} else {
		log.Println("category table already exists")
	}

	return nil
}
