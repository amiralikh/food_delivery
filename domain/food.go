package domain

type Food struct {
	ID            int64    `json:"id"`
	Name          string   `json:"name"`
	SupplierID    int64    `json:"supplier_id"`
	SupplierName  string   `json:"supplier_name"`
	CategoryID    int64    `json:"category_id"`
	CategoryName  string   `json:"category_name"`
	ImageUrl      string   `json:"image_url"`
	Description   string   `json:"description"`
	Price         int8     `json:"price"`
	DailyQuantity int8     `json:"daily_quantity"`
	Gallery       []*Image `json:"gallery"`
}

type Image struct {
	ID       int64  `json:"id"`
	FoodID   int64  `json:"food_id"`
	ImageURL string `json:"image_url"`
}
