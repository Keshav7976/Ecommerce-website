// controllers/category.go
package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/keshav7976/ecommerce/config"
	"github.com/keshav7976/ecommerce/models"
)

func GetCategories(w http.ResponseWriter, r *http.Request) {
	var categories []models.Category
	config.DB.Find(&categories)
	json.NewEncoder(w).Encode(categories)
}
