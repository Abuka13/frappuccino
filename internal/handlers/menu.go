package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"frappuccino/internal/db"
	"net/http"

	"github.com/lib/pq"
)

func CreateMenuItem(dbc *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		if r.Header.Get("Content-Type") != "application/json" {
			http.Error(w, "Content-Type must be application/json", http.StatusUnsupportedMediaType)
			return
		}

		var item db.MenuItem
		if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
			http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		// Validate required fields
		if item.Name == "" {
			http.Error(w, "name is required", http.StatusBadRequest)
			return
		}
		if item.Price <= 0 {
			http.Error(w, "price must be greater than 0", http.StatusBadRequest)
			return
		}
		if item.Category == "" {
			http.Error(w, "category is required", http.StatusBadRequest)
			return
		}

		query := `
			INSERT INTO menu_items (id,name, description, price, allergens, category, size)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
			RETURNING id
		`
		var id string
		err := dbc.QueryRowContext(r.Context(), query,
			item.ID,
			item.Name,
			item.Description,
			item.Price,
			pq.Array(item.Allergens),
			item.Category,
			item.Size,
		).Scan(&id)
		if err != nil {
			http.Error(w, "Failed to create menu item: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"id": id})
	}
}

func GetMenuItems(dbc *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		query := "SELECT id, name, description, price, allergens, category, size FROM menu_items"
		rows, err := dbc.QueryContext(r.Context(), query)
		if err != nil {
			http.Error(w, "Failed to fetch menu items", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var items []db.MenuItem
		for rows.Next() {
			var item db.MenuItem
			var allergens []string
			if err := rows.Scan(
				&item.ID,
				&item.Name,
				&item.Description,
				&item.Price,
				pq.Array(&allergens),
				&item.Category,
				&item.Size,
			); err != nil {
				http.Error(w, "Failed to scan menu item", http.StatusInternalServerError)
				return
			}
			item.Allergens = allergens
			items = append(items, item)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(items)
	}
}

func GetMenuItemByID(dbc *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		path := r.URL.Path
		var id string
		_, err := fmt.Sscanf(path, "/menu/%s", &id)
		if err != nil {
			http.Error(w, "Invalid menu item ID", http.StatusBadRequest)
			return
		}

		query := `
			SELECT id, name, description, price, allergens, category, size 
			FROM menu_items 
			WHERE id = $1
		`
		var item db.MenuItem
		var allergens []string
		err = dbc.QueryRowContext(r.Context(), query, id).Scan(
			&item.ID,
			&item.Name,
			&item.Description,
			&item.Price,
			pq.Array(&allergens),
			&item.Category,
			&item.Size,
		)
		if err == sql.ErrNoRows {
			http.Error(w, "Menu item not found", http.StatusNotFound)
			return
		} else if err != nil {
			http.Error(w, "Failed to fetch menu item", http.StatusInternalServerError)
			return
		}
		item.Allergens = allergens

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(item)
	}
}

func UpdateMenuItem(dbc *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		if r.Header.Get("Content-Type") != "application/json" {
			http.Error(w, "Content-Type must be application/json", http.StatusUnsupportedMediaType)
			return
		}

		path := r.URL.Path
		var id string
		_, err := fmt.Sscanf(path, "/menu/%s", &id)
		if err != nil {
			http.Error(w, "Invalid menu item ID", http.StatusBadRequest)
			return
		}

		var item db.MenuItem
		if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
			http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		// Validate required fields
		if item.Name == "" {
			http.Error(w, "name is required", http.StatusBadRequest)
			return
		}
		if item.Price <= 0 {
			http.Error(w, "price must be greater than 0", http.StatusBadRequest)
			return
		}
		if item.Category == "" {
			http.Error(w, "category is required", http.StatusBadRequest)
			return
		}

		query := `
			UPDATE menu_items 
			SET name = $2, description = $3, price = $4, allergens = $5, category = $6, size = $7
			WHERE id = $1
			RETURNING id
		`
		err = dbc.QueryRowContext(r.Context(), query,
			id,
			item.Name,
			item.Description,
			item.Price,
			pq.Array(item.Allergens),
			item.Category,
			item.Size,
		).Scan(&id)
		if err != nil {
			http.Error(w, "Failed to update menu item: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"id": id})
	}
}

func DeleteMenuItem(dbc *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		path := r.URL.Path
		var id string
		_, err := fmt.Sscanf(path, "/menu/%s", &id)
		if err != nil {
			http.Error(w, "Invalid menu item ID", http.StatusBadRequest)
			return
		}

		query := "DELETE FROM menu_items WHERE id = $1"
		result, err := dbc.ExecContext(r.Context(), query, id)
		if err != nil {
			http.Error(w, "Failed to delete menu item", http.StatusInternalServerError)
			return
		}

		rowsAffected, _ := result.RowsAffected()
		if rowsAffected == 0 {
			http.Error(w, "Menu item not found", http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

// func ToggleMenuItemAvailability(dbc *sql.DB) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		if r.Method != http.MethodPost {
// 			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
// 			return
// 		}

// 		path := r.URL.Path
// 		var id string
// 		_, err := fmt.Sscanf(path, "/menu_items/toggle/%s", &id)
// 		if err != nil {
// 			http.Error(w, "Invalid menu item ID", http.StatusBadRequest)
// 			return
// 		}

// 		query := `
// 			UPDATE menu_items 
// 			SET is_available = NOT is_available 
// 			WHERE id = $1
// 			RETURNING id, is_available
// 		`
// 		var (
// 			updatedID     string
// 			isAvailable bool
// 		)
// 		err = dbc.QueryRowContext(r.Context(), query, id).Scan(&updatedID, &isAvailable)
// 		if err == sql.ErrNoRows {
// 			http.Error(w, "Menu item not found", http.StatusNotFound)
// 			return
// 		} else if err != nil {
// 			http.Error(w, "Failed to toggle availability", http.StatusInternalServerError)
// 			return
// 		}

// 		w.Header().Set("Content-Type", "application/json")
// 		json.NewEncoder(w).Encode(map[string]interface{}{
// 			"id":           updatedID,
// 			"is_available": isAvailable,
// 		})
// 	}
// }