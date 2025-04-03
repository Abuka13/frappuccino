// Здесь реализуются обработчики для эндпоинтов, связанных с заказами.

package models

import (
    "database/sql"
    "encoding/json"
    "net/http"
    "frappuccino/internal/db"
)

func CreateOrder(dbc *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var order db.Order
        if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
            http.Error(w, "Invalid request body", http.StatusBadRequest)
            return
        }

        query := `INSERT INTO orders (customer_id, total_amount, status, payment_method) 
                  VALUES ($1, $2, $3, $4) RETURNING id`
        var orderID int
        err := dbc.QueryRow(query, order.CustomerID, order.TotalAmount, order.Status, order.PaymentMethod).Scan(&orderID)
        if err != nil {
            http.Error(w, "Failed to create order", http.StatusInternalServerError)
            return
        }

        w.WriteHeader(http.StatusCreated)
        json.NewEncoder(w).Encode(map[string]int{"order_id": orderID})
    }
}