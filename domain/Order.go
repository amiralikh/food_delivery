package domain

type Order struct {
	ID           int64        `json:"ID"`
	UserID       int64        `json:"user_id"`
	UserName     string       `json:"user_name"`
	SupplierID   int64        `json:"supplier_id"`
	SupplierName string       `json:"supplier_name"`
	AddressID    string       `json:"address_id"`
	TrackingID   string       `json:"tracking_id"`
	Status       string       `json:"status"`
	Price        float32      `json:"price"`
	CreatedAT    string       `json:"created_at"`
	Items        *[]OrderItem `json:"items"`
}

type OrderItem struct {
	ID          int64   `json:"id"`
	OrderID     int64   `json:"order_id"`
	FoodID      int64   `json:"food_id"`
	FoodName    string  `json:"food_name"`
	Quantity    int8    `json:"quantity"`
	SinglePrice float32 `json:"single_price"`
}
