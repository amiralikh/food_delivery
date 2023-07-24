package domain

type Supplier struct {
	ID           int64  `json:"ID"`
	Name         string `json:"name"`
	Address      string `json:"address"`
	Description  string `json:"description"`
	LogoUrl      string `json:"logo_url"`
	OpeningHour  string `json:"opening_hour"`
	ClosingHour  string `json:"closing_hour"`
	UserID       int64  `json:"user_id"`
	DeliveryTime string `json:"delivery_time"`
}
