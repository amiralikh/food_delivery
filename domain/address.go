package domain

type Address struct {
	ID      int64  `json:"id"`
	UserID  int64  `json:"user_id"`
	Name    string `json:"name"`
	Zip     string `json:"zip"`
	Phone   string `json:"phone"`
	Address string `json:"address"`
}
