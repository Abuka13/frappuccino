package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"frappuccino/internal/db"
	"net/http"
	
)

func CreateInventoryItem(dbc *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		if r.Header.Get("Content-Type") != "application/json" {
			http.Error(w, "Content-Type must be application/json", http.StatusUnsupportedMediaType)
			return
		}

		var item db.Inventory
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
		if item.Stock < 0 {
			http.Error(w, "stock cannot be negative", http.StatusBadRequest)
			return
		}
		if item.Price <= 0 {
			http.Error(w, "price must be greater than 0", http.StatusBadRequest)
			return
		}
		if item.UnitType == "" {
			http.Error(w, "unit_type is required", http.StatusBadRequest)
			return
		}

		query := `
			INSERT INTO inventory (name, stock, price, unit_type, last_updated)
			VALUES ($1, $2, $3, $4, NOW())
			RETURNING id
		`
		var id string
		err := dbc.QueryRowContext(r.Context(), query,
			item.Name,
			item.Stock,
			item.Price,
			item.UnitType,
		).Scan(&id)
		if err != nil {
			http.Error(w, "Failed to create inventory item: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"id": id})
	}
}

func GetInventoryItems(dbc *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		query := "SELECT id, name, stock, price, unit_type, last_updated FROM inventory"
		rows, err := dbc.QueryContext(r.Context(), query)
		if err != nil {
			http.Error(w, "Failed to fetch inventory items", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var items []db.Inventory
		for rows.Next() {
			var item db.Inventory
			if err := rows.Scan(
				&item.ID,
				&item.Name,
				&item.Stock,
				&item.Price,
				&item.UnitType,
				&item.LastUpdated,
			); err != nil {
				http.Error(w, "Failed to scan inventory item", http.StatusInternalServerError)
				return
			}
			items = append(items, item)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(items)
	}
}

func GetInventoryItemByID(dbc *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		path := r.URL.Path
		var id string
		_, err := fmt.Sscanf(path, "/inventory/%s", &id)
		if err != nil {
			http.Error(w, "Invalid inventory item ID", http.StatusBadRequest)
			return
		}

		query := `
			SELECT id, name, stock, price, unit_type, last_updated 
			FROM inventory 
			WHERE id = $1
		`
		var item db.Inventory
		err = dbc.QueryRowContext(r.Context(), query, id).Scan(
			&item.ID,
			&item.Name,
			&item.Stock,
			&item.Price,
			&item.UnitType,
			&item.LastUpdated,
		)
		if err == sql.ErrNoRows {
			http.Error(w, "Inventory item not found", http.StatusNotFound)
			return
		} else if err != nil {
			http.Error(w, "Failed to fetch inventory item", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(item)
	}
}

func UpdateInventoryItem(dbc *sql.DB) http.HandlerFunc {
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
		_, err := fmt.Sscanf(path, "/inventory/%s", &id)
		if err != nil {
			http.Error(w, "Invalid inventory item ID", http.StatusBadRequest)
			return
		}

		var item db.Inventory
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
		if item.Stock < 0 {
			http.Error(w, "stock cannot be negative", http.StatusBadRequest)
			return
		}
		if item.Price <= 0 {
			http.Error(w, "price must be greater than 0", http.StatusBadRequest)
			return
		}
		if item.UnitType == "" {
			http.Error(w, "unit_type is required", http.StatusBadRequest)
			return
		}

		query := `
			UPDATE inventory 
			SET name = $2, stock = $3, price = $4, unit_type = $5, last_updated = NOW()
			WHERE id = $1
			RETURNING id
		`
		err = dbc.QueryRowContext(r.Context(), query,
			id,
			item.Name,
			item.Stock,
			item.Price,
			item.UnitType,
		).Scan(&id)
		if err != nil {
			http.Error(w, "Failed to update inventory item: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"id": id})
	}
}

func DeleteInventoryItem(dbc *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		path := r.URL.Path
		var id string
		_, err := fmt.Sscanf(path, "/inventory/%s", &id)
		if err != nil {
			http.Error(w, "Invalid inventory item ID", http.StatusBadRequest)
			return
		}

		query := "DELETE FROM inventory WHERE id = $1"
		result, err := dbc.ExecContext(r.Context(), query, id)
		if err != nil {
			http.Error(w, "Failed to delete inventory item", http.StatusInternalServerError)
			return
		}

		rowsAffected, _ := result.RowsAffected()
		if rowsAffected == 0 {
			http.Error(w, "Inventory item not found", http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

// func RestockInventoryItem(dbc *sql.DB) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		if r.Method != http.MethodPost {
// 			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
// 			return
// 		}

// 		path := r.URL.Path
// 		var id string
// 		_, err := fmt.Sscanf(path, "/inventory/restock/%s", &id)
// 		if err != nil {
// 			http.Error(w, "Invalid inventory item ID", http.StatusBadRequest)
// 			return
// 		}

// 		type RestockRequest struct {
// 			Amount float64 `json:"amount"`
// 		}
// 		var req RestockRequest
// 		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// 			http.Error(w, "Invalid request body", http.StatusBadRequest)
// 			return
// 		}
// 		defer r.Body.Close()

// 		if req.Amount <= 0 {
// 			http.Error(w, "Amount must be greater than 0", http.StatusBadRequest)
// 			return
// 		}

// 		query := `
// 			UPDATE inventory 
// 			SET stock = stock + $2, last_updated = NOW()
// 			WHERE id = $1
// 			RETURNING id, stock
// 		`
// 		var (
// 			updatedID string
// 			newStock  float64
// 		)
// 		err = dbc.QueryRowContext(r.Context(), query, id, req.Amount).Scan(&updatedID, &newStock)
// 		if err == sql.ErrNoRows {
// 			http.Error(w, "Inventory item not found", http.StatusNotFound)
// 			return
// 		} else if err != nil {
// 			http.Error(w, "Failed to restock item", http.StatusInternalServerError)
// 			return
// 		}

// 		w.Header().Set("Content-Type", "application/json")
// 		json.NewEncoder(w).Encode(map[string]interface{}{
// 			"id":    updatedID,
// 			"stock": newStock,
// 		})
// 	}
// }