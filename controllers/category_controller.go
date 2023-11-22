package controllers

import (
	"encoding/json"
	"net/http"
	"time"

	"strconv"

	"github.com/Hafidzurr/project4_group2_glng-ks-08/config"
	"github.com/Hafidzurr/project4_group2_glng-ks-08/models"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

// CreateCategory - Create a new category
func CreateCategory(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Authenticate and authorize the user
		user, err := config.ExtractUserFromToken(r, db)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		if user.Role != "admin" {
			http.Error(w, "Unauthorized access", http.StatusForbidden)
			return
		}

		var requestBody struct {
			Type string `json:"type"`
		}
		err = json.NewDecoder(r.Body).Decode(&requestBody)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		category := models.Category{
			Type:              requestBody.Type,
			SoldProductAmount: 0,
			CreatedAt:         time.Now(),
		}

		result := db.Create(&category)
		if result.Error != nil {
			http.Error(w, result.Error.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		config.SendJSONResponse(w, map[string]interface{}{
			"id":                  category.ID,
			"type":                category.Type,
			"sold_product_amount": category.SoldProductAmount,
			"created_at":          category.CreatedAt.Format(time.RFC3339),
		})
	}
}

// GetCategories
func GetCategories(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Authenticate and authorize the user
		user, err := config.ExtractUserFromToken(r, db)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		if user.Role != "admin" {
			http.Error(w, "Unauthorized access", http.StatusForbidden)
			return
		}

		var categories []models.Category
		result := db.Preload("Products").Find(&categories)
		if result.Error != nil {
			http.Error(w, result.Error.Error(), http.StatusInternalServerError)
			return
		}

		// Create a structured response with ordered fields
		var responseCategories []map[string]interface{}
		for _, category := range categories {
			categoryData := map[string]interface{}{
				"id":                  category.ID,
				"type":                category.Type,
				"sold_product_amount": category.SoldProductAmount,
				"created_at":          category.CreatedAt.Format(time.RFC3339),
				"updated_at":          category.UpdatedAt.Format(time.RFC3339),
				"products":            []map[string]interface{}{},
			}

			// Add product details to the response with ordered fields
			for _, product := range category.Products {
				productData := map[string]interface{}{
					"id":         product.ID,
					"title":      product.Title,
					"price":      product.Price,
					"stock":      product.Stock,
					"created_at": product.CreatedAt.Format(time.RFC3339),
					"updated_at": product.UpdatedAt.Format(time.RFC3339),
				}
				categoryData["products"] = append(categoryData["products"].([]map[string]interface{}), productData)
			}

			responseCategories = append(responseCategories, categoryData)
		}

		// Send the JSON response with ordered fields
		config.SendJSONResponse(w, responseCategories)
	}
}

// UpdateCategory - Update category by ID
func UpdateCategory(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Authenticate and authorize the user
		user, err := config.ExtractUserFromToken(r, db)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		if user.Role != "admin" {
			http.Error(w, "Unauthorized access", http.StatusForbidden)
			return
		}

		categoryID, err := strconv.Atoi(mux.Vars(r)["categoryId"])
		if err != nil {
			http.Error(w, "Invalid category ID", http.StatusBadRequest)
			return
		}

		var requestBody struct {
			Type string `json:"type"`
		}
		err = json.NewDecoder(r.Body).Decode(&requestBody)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var category models.Category
		result := db.First(&category, categoryID)
		if result.Error != nil {
			http.Error(w, "Category not found", http.StatusNotFound)
			return
		}

		category.Type = requestBody.Type
		category.UpdatedAt = time.Now()

		db.Save(&category)

		config.SendJSONResponse(w, map[string]interface{}{
			"id":                  category.ID,
			"type":                category.Type,
			"sold_product_amount": category.SoldProductAmount,
			"updated_at":          category.UpdatedAt.Format(time.RFC3339),
		})
	}
}

// DeleteCategory - Delete category by ID
func DeleteCategory(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Authenticate and authorize the user
		user, err := config.ExtractUserFromToken(r, db)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		if user.Role != "admin" {
			http.Error(w, "Unauthorized access", http.StatusForbidden)
			return
		}

		categoryID, err := strconv.Atoi(mux.Vars(r)["categoryId"])
		if err != nil {
			http.Error(w, "Invalid category ID", http.StatusBadRequest)
			return
		}

		var category models.Category
		result := db.First(&category, categoryID)
		if result.Error != nil {
			http.Error(w, "Category not found", http.StatusNotFound)
			return
		}

		db.Delete(&category)

		config.SendJSONResponse(w, map[string]interface{}{
			"message": "Category has been successfully deleted",
		})
	}
}
