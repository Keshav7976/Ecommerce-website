// File: controllers/item.go
package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/keshav7976/ecommerce/config"
	"github.com/keshav7976/ecommerce/models"
	"gorm.io/gorm"
)

// Helper function to find a category by name
func getCategoryByName(name string) (*models.Category, error) {
	var category models.Category
	if err := config.DB.Where("name = ?", name).First(&category).Error; err != nil {
		return nil, err
	}
	return &category, nil
}

// GET /items
func GetItems(w http.ResponseWriter, r *http.Request) {
	var items []models.Item
	query := config.DB

	// filters
	categoryID := r.URL.Query().Get("category_id")
	if categoryID != "" {
		query = query.Where("category_id = ?", categoryID)
	}

	minPrice := r.URL.Query().Get("minPrice")
	maxPrice := r.URL.Query().Get("maxPrice")
	if minPrice != "" && maxPrice != "" {
		min, err := strconv.ParseFloat(minPrice, 64)
		if err != nil {
			http.Error(w, `{"error":"Invalid minPrice"}`, http.StatusBadRequest)
			return
		}
		max, err := strconv.ParseFloat(maxPrice, 64)
		if err != nil {
			http.Error(w, `{"error":"Invalid maxPrice"}`, http.StatusBadRequest)
			return
		}
		query = query.Where("price BETWEEN ? AND ?", min, max)
	}

	if err := query.Preload("Category").Find(&items).Error; err != nil {
		http.Error(w, `{"error":"Failed to fetch items"}`, http.StatusInternalServerError)
		return
	}

	var cleanItems []models.ItemResponse
	for _, item := range items {
		categoryName := ""
		if item.Category != nil {
			categoryName = item.Category.Name
		}
		cleanItems = append(cleanItems, models.ItemResponse{
			ID:       item.ID,
			Name:     item.Name,
			Price:    item.Price,
			ImageURL: item.ImageURL,
			Category: categoryName,
		})
	}

	json.NewEncoder(w).Encode(cleanItems)
}

// POST /items (any authenticated user)
func CreateItem(w http.ResponseWriter, r *http.Request) {
	_, ok := r.Context().Value("user_id").(uint)
	if !ok {
		http.Error(w, `{"error":"Unauthorized"}`, http.StatusForbidden)
		return
	}

	var input struct {
		Name         string  `json:"name"`
		Price        float64 `json:"price"`
		ImageURL     string  `json:"image_url"`
		CategoryName string  `json:"category_name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, `{"error":"Invalid input"}`, http.StatusBadRequest)
		return
	}

	category, err := getCategoryByName(input.CategoryName)
	if err != nil {
		http.Error(w, `{"error":"Category not found"}`, http.StatusBadRequest)
		return
	}

	item := models.Item{
		Name:       input.Name,
		Price:      input.Price,
		ImageURL:   input.ImageURL,
		CategoryID: category.ID,
	}

	if result := config.DB.Create(&item); result.Error != nil {
		http.Error(w, `{"error":"Failed to create item"}`, http.StatusInternalServerError)
		return
	}

	response := models.ItemResponse{
		ID:       item.ID,
		Name:     item.Name,
		Price:    item.Price,
		ImageURL: item.ImageURL,
		Category: category.Name,
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// PUT /items/{id} (any authenticated user)
func UpdateItem(w http.ResponseWriter, r *http.Request) {
	_, ok := r.Context().Value("user_id").(uint)
	if !ok {
		http.Error(w, `{"error":"Unauthorized"}`, http.StatusForbidden)
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]

	var input struct {
		Name         string  `json:"name"`
		Price        float64 `json:"price"`
		ImageURL     string  `json:"image_url"`
		CategoryName string  `json:"category_name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, `{"error":"Invalid input"}`, http.StatusBadRequest)
		return
	}

	var item models.Item
	if err := config.DB.First(&item, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			http.Error(w, `{"error":"Item not found"}`, http.StatusNotFound)
			return
		}
		http.Error(w, `{"error":"Failed to find item"}`, http.StatusInternalServerError)
		return
	}

	if input.Name != "" {
		item.Name = input.Name
	}
	if input.Price != 0 {
		item.Price = input.Price
	}
	if input.ImageURL != "" {
		item.ImageURL = input.ImageURL
	}
	if input.CategoryName != "" {
		category, err := getCategoryByName(input.CategoryName)
		if err != nil {
			http.Error(w, `{"error":"Category not found"}`, http.StatusBadRequest)
			return
		}
		item.CategoryID = category.ID
	}

	if result := config.DB.Save(&item); result.Error != nil {
		http.Error(w, `{"error":"Failed to update item"}`, http.StatusInternalServerError)
		return
	}

	var category models.Category
	config.DB.First(&category, item.CategoryID)

	response := models.ItemResponse{
		ID:       item.ID,
		Name:     item.Name,
		Price:    item.Price,
		ImageURL: item.ImageURL,
		Category: category.Name,
	}

	json.NewEncoder(w).Encode(response)
}

// DELETE /items/{id} (any authenticated user)
func DeleteItem(w http.ResponseWriter, r *http.Request) {
	_, ok := r.Context().Value("user_id").(uint)
	if !ok {
		http.Error(w, `{"error":"Unauthorized"}`, http.StatusForbidden)
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]

	if err := config.DB.Delete(&models.Item{}, id).Error; err != nil {
		http.Error(w, `{"error":"Failed to delete item"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
