package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"frappuccino/internal/db"
	"log"
	"net/http"
)

func UpdateOrderByID(dbc *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // Check method first
        if r.Method != http.MethodPut {
            http.Error(w, "Invalid Method Request", http.StatusMethodNotAllowed)
            return
        }

        // Check content type
        contentType := r.Header.Get("Content-Type")
        if contentType != "application/json" {
            http.Error(w, "Content-Type must be application/json", http.StatusUnsupportedMediaType)
            return
        }

        // Extract order ID from URL
        path := r.URL.Path
        var orderID string
        _, err := fmt.Sscanf(path, "/orders/%s", &orderID)
        if err != nil {
            http.Error(w, "Invalid order ID", http.StatusBadRequest)
            return
        }

        // Parse request body
        var order db.Order
        if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
            http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
            return
        }
        defer r.Body.Close()

        // Validate required fields
        if order.CustomerID == 0 {
            http.Error(w, "customer_id is required", http.StatusBadRequest)
            return
        }
        if order.TotalAmount <= 0 {
            http.Error(w, "total_amount must be greater than 0", http.StatusBadRequest)
            return
        }
        if order.PaymentMethod == "" {
            http.Error(w, "payment_method is required", http.StatusBadRequest)
            return
        }

        // Set default status if not provided
        if order.Status == "" {
            order.Status = "open"
        }

        // Prepare special instructions
        var specialInstructions []byte
        if order.SpecialInstructions != nil {
            specialInstructions, err = json.Marshal(order.SpecialInstructions)
            if err != nil {
                http.Error(w, "Invalid special_instructions format", http.StatusBadRequest)
                return
            }
        }

        // Execute update query
        query := `UPDATE orders 
                 SET customer_id = $2, total_amount = $3, status = $4, 
                     special_instructions = $5, payment_method = $6, updated_at = NOW()
                 WHERE id = $1
                 RETURNING id`
        
        err = dbc.QueryRowContext(r.Context(), query,
            orderID, // $1
            order.CustomerID, // $2
            order.TotalAmount, // $3
            order.Status, // $4
            specialInstructions, // $5 (JSONB)
            order.PaymentMethod, // $6
        ).Scan(&orderID)
        
        if err != nil {
            http.Error(w, "Failed to update order: "+err.Error(), http.StatusInternalServerError)
            return
        }

        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)
        json.NewEncoder(w).Encode(map[string]string{"order_id": orderID})
    }
}

func DeleteOrder(dbc *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		var OrderID string
		_,err := fmt.Sscanf(path,"/orders/%s", &OrderID)
		if err != nil {
			http.Error(w, "Invalid order ID", http.StatusBadRequest)
            return
		}
		if r.Method != http.MethodDelete {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		// Verify content type
		if r.Header.Get("Content-Type") != "application/json" {
			http.Error(w, "Content-Type must be application/json", http.StatusUnsupportedMediaType)
			return
		}
		
		defer r.Body.Close()
		
		query := `DELETE FROM orders WHERE id = $1`
		result, err := dbc.ExecContext(r.Context(), query, OrderID)
		if err != nil {
			http.Error(w, "Failed to delete order: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Check if any rows were affected
		rowsAffected, err := result.RowsAffected()
		if err != nil {
			http.Error(w, "Failed to check deletion: "+err.Error(), http.StatusInternalServerError)
			return
		}

		if rowsAffected == 0 {
			http.Error(w, "Order not found", http.StatusNotFound)
			return
		}

		// Return success response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"message":  "Order deleted successfully",
			"order_id": OrderID,
		})
	}
}

func CreateOrder(dbc *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		// Verify content type
		if r.Header.Get("Content-Type") != "application/json" {
			http.Error(w, "Content-Type must be application/json", http.StatusUnsupportedMediaType)
			return
		}
		var order db.Order
		if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
			http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
			return
		}
		defer r.Body.Close()
		// Validate required fields
		if order.CustomerID == 0 {
			http.Error(w, "customer_id is required", http.StatusBadRequest)
			return
		}
		if order.TotalAmount <= 0 {
			http.Error(w, "total_amount must be greater than 0", http.StatusBadRequest)
			return
		}
		if order.PaymentMethod == "" {
			http.Error(w, "payment_method is required", http.StatusBadRequest)
			return
		}
		// Set default status if not provided
		if order.Status == "" {
			order.Status = "open"
		}
		// Insert into database
		query := `
            INSERT INTO orders (customer_id, total_amount, status, special_instructions, payment_method)
            VALUES ($1, $2, $3, $4, $5)
            RETURNING id
        `
		var orderID int
		err := dbc.QueryRowContext(r.Context(), query,
			order.CustomerID,
			order.TotalAmount,
			order.Status,
			order.SpecialInstructions, // This will be stored as JSONB in PostgreSQL
			order.PaymentMethod,
		).Scan(&orderID)
		if err != nil {
			http.Error(w, "Failed to create order: "+err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]int{"order_id": orderID})
	}
}


func GetOrders(dbс *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := "SELECT id, customer_id, total_amount, status, payment_method, created_at, updated_at FROM orders"
		rows, err := dbс.Query(query)
		if err != nil {
			http.Error(w, "Failed to fetch orders", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var orders []db.Order
		for rows.Next() {
			var order db.Order
			if err := rows.Scan(
				&order.ID,
				&order.CustomerID,
				&order.TotalAmount,
				&order.Status,
				&order.PaymentMethod,
				&order.CreatedAt,
				&order.UpdatedAt,
			); err != nil {
				http.Error(w, "Failed to scan order", http.StatusInternalServerError)
				return
			}
			orders = append(orders, order)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(orders)
	}
}

func GetOrderByID(dbс *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract the "id" parameter from the URL path
		path := r.URL.Path
		var orderID string
		_, err := fmt.Sscanf(path, "/orders/%s", &orderID)
		if err != nil {
			http.Error(w, "Invalid order ID", http.StatusBadRequest)
			return
		}

		query := `
            SELECT id, customer_id, total_amount, status, payment_method, created_at, updated_at 
            FROM orders 
            WHERE id = $1
        `
		var order db.Order
		err = dbс.QueryRow(query, orderID).Scan(
			&order.ID,
			&order.CustomerID,
			&order.TotalAmount,
			&order.Status,
			&order.PaymentMethod,
			&order.CreatedAt,
			&order.UpdatedAt,
		)
		if err == sql.ErrNoRows {
			http.Error(w, "Order not found", http.StatusNotFound)
			return
		} else if err != nil {
			http.Error(w, "Failed to fetch order", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(order)
	}
}

func CloseOrder(dbс *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract the "id" parameter from the URL path
		path := r.URL.Path
		var orderID string
		_, err := fmt.Sscanf(path, "/orders/close/%s", &orderID)
		if err != nil {
			http.Error(w, "Invalid order ID", http.StatusBadRequest)
			log.Println(err)
			return
		}

		// Обновляем статус заказа на "closed"
		query := `
            UPDATE orders 
            SET status = 'closed', updated_at = NOW() 
            WHERE id = $1 AND status = 'open'
        `
		result, err := dbс.Exec(query, orderID)
		if err != nil {
			http.Error(w, "Failed to close order", http.StatusInternalServerError)
			return
		}

		rowsAffected, _ := result.RowsAffected()
		if rowsAffected == 0 {
			http.Error(w, "Order not found or already closed", http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "Order closed successfully"})
	}
}

func GetNumberOfOrderedItems(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		startDate := r.URL.Query().Get("startDate")
		endDate := r.URL.Query().Get("endDate")

		query := `
            SELECT 
                order_items.menu_item_id, 
                SUM(order_items.quantity) AS total_quantity,  
                orders.created_at
            FROM 
                order_items
            JOIN 
                orders ON order_items.order_id = orders.id
            WHERE 
                orders.created_at BETWEEN $1 AND $2
            GROUP BY 
                order_items.menu_item_id, orders.created_at;

        `
		rows, err := db.Query(query, startDate, endDate)
		if err != nil {
			http.Error(w, "Failed to fetch ordered items", http.StatusInternalServerError)
			log.Println(err)
			return
		}
		defer rows.Close()

		type OrderItem struct {
			ID        string `json:"menuItemID"`
			Total     int    `json:"totalQuantity"`
			CreatedAt string `json:"createdAt"`
		}
		var items []OrderItem
		for rows.Next() {
			var menuItemID string
			var totalQuantity int
			var createdAt string

			if err := rows.Scan(&menuItemID, &totalQuantity, &createdAt); err != nil {
				http.Error(w, "Failed to scan ordered items", http.StatusInternalServerError)
				log.Println(err)
				return
			}
			item := OrderItem{
				ID:        menuItemID,
				Total:     totalQuantity,
				CreatedAt: createdAt,
			}
			items = append(items, item)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(items)
	}
}

// PUT DELETE
