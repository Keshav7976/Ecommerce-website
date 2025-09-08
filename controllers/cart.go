package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/keshav7976/ecommerce/config"
	"github.com/keshav7976/ecommerce/models"
)

// --- Add item to cart ---
func AddToCart(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(uint)
	if !ok {
		http.Error(w, `{"error":"User ID not found in context"}`, http.StatusUnauthorized)
		return
	}

	var req struct {
		ItemID uint `json:"item_id"`
		Qty    int  `json:"qty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"Invalid request body"}`, http.StatusBadRequest)
		return
	}

	if req.ItemID == 0 || req.Qty <= 0 {
		http.Error(w, `{"error":"Invalid item ID or quantity"}`, http.StatusBadRequest)
		return
	}

	var existing models.CartItem
	err := config.DB.Where("user_id = ? AND item_id = ?", userID, req.ItemID).First(&existing).Error
	if err == nil {
		// item exists, update qty
		existing.Qty += req.Qty
		config.DB.Save(&existing)
	} else {
		// create new
		newItem := models.CartItem{UserID: userID, ItemID: req.ItemID, Qty: req.Qty}
		if err := config.DB.Create(&newItem).Error; err != nil {
			http.Error(w, `{"error":"Failed to add item to cart"}`, http.StatusInternalServerError)
			return
		}
		existing = newItem
	}

	// Manually load item info
	var item models.Item
	if err := config.DB.First(&item, existing.ItemID).Error; err != nil {
		http.Error(w, `{"error":"Item not found"}`, http.StatusInternalServerError)
		return
	}

	resp := map[string]interface{}{
		"id":        existing.ID,
		"qty":       existing.Qty,
		"item_id":   item.ID,
		"name":      item.Name,
		"price":     item.Price,
		"image_url": item.ImageURL,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// --- Get cart items ---
func GetCart(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(uint)
	if !ok {
		http.Error(w, `{"error":"User ID not found in context"}`, http.StatusUnauthorized)
		return
	}

	// Log user ID
	fmt.Printf("Fetching cart for userID: %d\n", userID)

	var cart []models.CartItem
	// Preload Item and its Category
	err := config.DB.Preload("Item.Category").Where("user_id = ?", userID).Find(&cart).Error
	if err != nil {
		fmt.Printf("DB error: %v\n", err)
		http.Error(w, `{"error":"Failed to fetch cart"}`, http.StatusInternalServerError)
		return
	}

	fmt.Printf("Number of cart items fetched: %d\n", len(cart))
	for i, c := range cart {
		fmt.Printf("CartItem %d -> ID:%d, UserID:%d, ItemID:%d, Qty:%d\n", i, c.ID, c.UserID, c.ItemID, c.Qty)
		if c.Item.ID == 0 {
			fmt.Println("  -> Associated item is empty")
		} else {
			fmt.Printf("  -> Item: ID:%d, Name:%s, Price:%.2f\n", c.Item.ID, c.Item.Name, c.Item.Price)
		}
	}

	// Respond
	type CartResp struct {
		ID       uint    `json:"id"`
		Qty      int     `json:"qty"`
		ItemID   uint    `json:"item_id"`
		Name     string  `json:"name"`
		Price    float64 `json:"price"`
		ImageURL string  `json:"image_url"`
		Category string  `json:"category"`
	}

	resp := []CartResp{}
	for _, c := range cart {
		if c.Item.ID == 0 {
			continue
		}
		categoryName := ""
		if c.Item.Category != nil {
			categoryName = c.Item.Category.Name
		}
		resp = append(resp, CartResp{
			ID:       c.ID,
			Qty:      c.Qty,
			ItemID:   c.Item.ID,
			Name:     c.Item.Name,
			Price:    c.Item.Price,
			ImageURL: c.Item.ImageURL,
			Category: categoryName,
		})
	}

	fmt.Printf("Number of cart items sent to frontend: %d\n", len(resp))
	json.NewEncoder(w).Encode(resp)
}

// --- Remove an item from cart ---
func RemoveFromCart(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(uint)
	if !ok {
		http.Error(w, `{"error":"User ID not found in context"}`, http.StatusUnauthorized)
		return
	}

	itemIDStr := r.URL.Query().Get("item_id")
	if itemIDStr == "" {
		http.Error(w, `{"error":"Item ID is required"}`, http.StatusBadRequest)
		return
	}

	itemID, err := strconv.ParseUint(itemIDStr, 10, 32)
	if err != nil {
		http.Error(w, `{"error":"Invalid item ID"}`, http.StatusBadRequest)
		return
	}

	if err := config.DB.Where("user_id = ? AND item_id = ?", userID, uint(itemID)).Delete(&models.CartItem{}).Error; err != nil {
		http.Error(w, `{"error":"Failed to remove item from cart"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Item removed from cart"})
}
