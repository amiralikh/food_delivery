package repository

import (
	"database/sql"
	"errors"
	"foodDelivery/domain"
	"time"
)

type OrderRepository interface {
}

type orderRepository struct {
	db *sql.DB
}

var (
	ErrOrderNotFound = errors.New("order not found")
)

func NewOrderRepository(db *sql.DB) OrderRepository {
	return &orderRepository{
		db: db,
	}
}

func (or *orderRepository) GetUserOrders(userId int64) (*[]domain.Order, error) {
	orders := []domain.Order{}

	query := `
		SELECT o.id, o.user_id, u.name AS user_name, o.supplier_id, s.name AS supplier_name,
			o.tracking_id, o.status, o.price, o.created_at
		FROM orders o
		INNER JOIN users u ON o.user_id = u.id
		INNER JOIN suppliers s ON o.supplier_id = s.id
		WHERE o.user_id = $1
		ORDER BY o.created_at DESC
	`

	rows, err := or.db.Query(query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var order domain.Order
		err := rows.Scan(
			&order.ID,
			&order.UserID,
			&order.UserName,
			&order.SupplierID,
			&order.SupplierName,
			&order.TrackingID,
			&order.Status,
			&order.Price,
			&order.CreatedAT,
		)
		if err != nil {
			return nil, err
		}

		orders = append(orders, order)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &orders, nil
}

func (or *orderRepository) GetOrderWithItems(orderID int64) (*domain.Order, error) {
	order := &domain.Order{}

	orderQuery := `
		SELECT o.id, o.user_id, u.name AS user_name, o.supplier_id, s.name AS supplier_name,
			o.tracking_id, o.status, o.price, o.created_at
		FROM orders o
		INNER JOIN users u ON o.user_id = u.id
		INNER JOIN suppliers s ON o.supplier_id = s.id
		WHERE o.id = $1
	`
	err := or.db.QueryRow(orderQuery, orderID).Scan(
		&order.ID,
		&order.UserID,
		&order.UserName,
		&order.SupplierID,
		&order.SupplierName,
		&order.TrackingID,
		&order.Status,
		&order.Price,
		&order.CreatedAT,
	)
	if err != nil {
		return nil, err
	}

	itemsQuery := `
		SELECT oi.id, oi.order_id, oi.food_id, f.name AS food_name, oi.quantity, oi.single_price
		FROM order_items oi
		INNER JOIN foods f ON oi.food_id = f.id
		WHERE oi.order_id = $1
	`
	rows, err := or.db.Query(itemsQuery, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orderItems []domain.OrderItem
	for rows.Next() {
		var item domain.OrderItem
		err := rows.Scan(
			&item.ID,
			&item.OrderID,
			&item.FoodID,
			&item.FoodName,
			&item.Quantity,
			&item.SinglePrice,
		)
		if err != nil {
			return nil, err
		}

		orderItems = append(orderItems, item)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	order.Items = &orderItems
	return order, nil
}

func (or *orderRepository) SubmitOrder(order *domain.Order) error {
	tx, err := or.db.Begin()
	if err != nil {
		return err
	}
	order.CreatedAT = time.Now().UTC().Format("2006-01-02 15:04:05")
	order.Price = 0.0
	for _, item := range *order.Items {
		order.Price += item.SinglePrice * float32(item.Quantity)
	}

	orderQuery := `
		INSERT INTO orders (user_id, supplier_id, tracking_id, status, price, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`
	var orderID int64
	err = tx.QueryRow(orderQuery, order.UserID, order.SupplierID, order.TrackingID, order.Status, order.Price, order.CreatedAT).Scan(&orderID)
	if err != nil {
		tx.Rollback()
		return err
	}

	itemQuery := `
		INSERT INTO order_items (order_id, food_id, quantity, single_price)
		VALUES ($1, $2, $3, $4)
	`
	for _, item := range *order.Items {
		_, err = tx.Exec(itemQuery, orderID, item.FoodID, item.Quantity, item.SinglePrice)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}