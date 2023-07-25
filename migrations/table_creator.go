package migrations

import (
	"database/sql"
	"fmt"
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
	// Check if the categories table exists.
	tableExists := false
	err := db.QueryRow("SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'categories')").Scan(&tableExists)
	if err != nil {
		return err
	}

	// Create the categories table if it doesn't exist.
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

func CreateSuppliersTable(db *sql.DB) error {
	tableExists := false
	err := db.QueryRow("SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'suppliers')").Scan(&tableExists)
	if err != nil {
		return err
	}

	// Create the suppliers table if it doesn't exist.
	if !tableExists {
		createUsersTable := `
			CREATE TABLE suppliers (
				id SERIAL PRIMARY KEY,
				name VARCHAR(255) NOT NULL,
				address VARCHAR(255) NOT NULL,
				description VARCHAR(512) NOT NULL,
				logo_url VARCHAR(255) NOT NULL,
				opening_hour VARCHAR(50) NOT NULL,
				closing_hour VARCHAR(50) NOT NULL,
            	user_id INTEGER NOT NULL REFERENCES users(id),
				delivery_time VARCHAR(50) NOT NULL
			)
		`

		_, err = db.Exec(createUsersTable)
		if err != nil {
			return err
		}

		log.Println("suppliers table created successfully")
	} else {
		log.Println("suppliers table already exists")
	}

	return nil
}

func CreateFoodsTable(db *sql.DB) error {
	foodTableExists := false
	err := db.QueryRow("SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'foods')").Scan(&foodTableExists)
	if err != nil {
		return err
	}
	if !foodTableExists {
		foodsTableQuery := `
		CREATE TABLE IF NOT EXISTS foods (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			supplier_id INT NOT NULL REFERENCES suppliers(id),
			category_id INT NOT NULL REFERENCES categories(id),
			image_url VARCHAR(255) NOT NULL,
			description TEXT NOT NULL,
			price INT NOT NULL,
			daily_quantity INT NOT NULL
		)
	`
		_, err = db.Exec(foodsTableQuery)
		if err != nil {
			return err
		}
		log.Println("foods table created successfully")
	} else {
		log.Println("foods table already exists")
	}
	return nil
}

func CreateGalleryTable(db *sql.DB) error {
	galleryTableExists := false
	err := db.QueryRow("SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'gallery')").Scan(&galleryTableExists)
	if err != nil {
		return err
	}
	if !galleryTableExists {
		galleryTableQuery := `
		CREATE TABLE IF NOT EXISTS gallery (
			id SERIAL PRIMARY KEY,
			food_id INT NOT NULL REFERENCES foods(id),
			image_url VARCHAR(255) NOT NULL
		)
	`
		_, err = db.Exec(galleryTableQuery)
		if err != nil {
			return fmt.Errorf("failed to create gallery table: %v", err)
		}
		log.Println("gallery table created successfully")
	} else {
		log.Println("gallery table already exists")
	}
	return nil
}

func CreateAddressesTable(db *sql.DB) error {
	ordersTableExists := false
	err := db.QueryRow("SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'addresses')").Scan(&ordersTableExists)
	if err != nil {
		return err
	}
	if !ordersTableExists {
		ordersTableQuery := `
		CREATE TABLE IF NOT EXISTS addresses (
			id SERIAL PRIMARY KEY,
			user_id BIGINT NOT NULL REFERENCES users(id),
			name VARCHAR(50) NOT NULL,
			zip VARCHAR(50) NOT NULL,
			phone VARCHAR(50) NOT NULL,
			Address VARCHAR(225) NOT NULL
		)
	`
		_, err = db.Exec(ordersTableQuery)
		if err != nil {
			return fmt.Errorf("failed to create orders table: %v", err)
		}
		log.Println("addresses table created successfully")
	} else {
		log.Println("addresses table already exists")
	}
	return nil
}

func CreateOrdersTable(db *sql.DB) error {
	ordersTableExists := false
	err := db.QueryRow("SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'orders')").Scan(&ordersTableExists)
	if err != nil {
		return err
	}
	if !ordersTableExists {
		ordersTableQuery := `
		CREATE TABLE IF NOT EXISTS orders (
			id SERIAL PRIMARY KEY,
			user_id BIGINT NOT NULL REFERENCES users(id),
			supplier_id BIGINT NOT NULL REFERENCES suppliers(id),
		    address_id BIGINT NOT NULL REFERENCES addresses(id),
			tracking_id VARCHAR(50) NOT NULL,
			status VARCHAR(50) NOT NULL,
			price NUMERIC(10, 2) NOT NULL,
			created_at TIMESTAMP NOT NULL
		)
	`
		_, err = db.Exec(ordersTableQuery)
		if err != nil {
			return fmt.Errorf("failed to create orders table: %v", err)
		}
		log.Println("orders table created successfully")
	} else {
		log.Println("orders table already exists")
	}
	return nil
}

func CreateOrderItemsTable(db *sql.DB) error {
	orderItemsTableExists := false
	err := db.QueryRow("SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'order_items')").Scan(&orderItemsTableExists)
	if err != nil {
		return err
	}
	if !orderItemsTableExists {
		orderItemsTableQuery := `
		CREATE TABLE IF NOT EXISTS order_items (
			id SERIAL PRIMARY KEY,
			order_id BIGINT NOT NULL REFERENCES orders(id),
			food_id BIGINT NOT NULL REFERENCES foods(id),
			quantity SMALLINT NOT NULL,
			single_price NUMERIC(10, 2) NOT NULL
		)
	`
		_, err = db.Exec(orderItemsTableQuery)
		if err != nil {
			return fmt.Errorf("failed to create order_items table: %v", err)
		}
		log.Println("order_items table created successfully")
	} else {
		log.Println("order_items table already exists")
	}
	return nil
}
