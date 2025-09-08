// models/user.go
package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `gorm:"unique;size:255;not null"`
	Password string `gorm:"not null"`
	IsAdmin  bool   `gorm:"default:false"`
}
