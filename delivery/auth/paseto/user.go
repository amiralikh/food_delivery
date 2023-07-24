package paseto

type User struct {
	ID     int64  `json:"id"`
	Email  string `json:"email"`
	Status string `json:"status"`
}

type UserService interface {
	Authenticate(email, password string) (*User, error)
	GetUserByID(userID int64) (*User, error)
}
