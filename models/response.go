// models/response.go
package models

type ItemResponse struct {
	ID       uint    `json:"id"`
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	ImageURL string  `json:"image_url"`
	Category string  `json:"category"`
}
