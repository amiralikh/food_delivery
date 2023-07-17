package domain

type Category struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	ImageUrl string `json:"image_url"`
}
