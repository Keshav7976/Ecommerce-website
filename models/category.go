// models/category.go
package models

import "gorm.io/gorm"

type Category struct {
	gorm.Model
	Name  string `gorm:"unique;size:255;not null"`
	Items []Item
}
