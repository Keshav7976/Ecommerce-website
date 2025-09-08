// File: models/item.go
package models

import "gorm.io/gorm"

type Item struct {
	gorm.Model
	Name       string    `json:"name"`
	Price      float64   `json:"price"`
	ImageURL   string    `json:"image_url"`
	CategoryID uint      `json:"-"`
	Category   *Category `json:"category,omitempty" gorm:"foreignKey:CategoryID"`
}