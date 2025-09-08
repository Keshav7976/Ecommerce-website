package models

import "gorm.io/gorm"

type CartItem struct {
    gorm.Model
    UserID uint
    ItemID uint
    Qty    int
    Item   *Item `gorm:"foreignKey:ItemID;references:ID"`
}

