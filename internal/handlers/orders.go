package handlers
import (
    "frappuccino/internal/db"
    "database/sql"
    "encoding/json"
    "fmt"
    "net/http"
    
)

func CreateOrder(dbс *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var order db.Order
        if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
            http.Error(w, "Invalid request body", http.StatusBadRequest)
            return
        }

        // Вставка заказа в базу данных
        query := `
            INSERT INTO orders (customer_id, total_amount, status, special_instructions, payment_method)
            VALUES ($1, $2, $3, $4, $5)
            RETURNING id
        `
        var orderID int
        err := dbс.QueryRow(
            query,
            order.CustomerID,
            order.TotalAmount,
            order.Status,
            order.SpecialInstructions,
            order.PaymentMethod,
        ).Scan(&orderID)
        if err != nil {
            http.Error(w, "Failed to create order", http.StatusInternalServerError)
            return
        }

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
        _, err := fmt.Sscanf(path, "/orders/%s/close", &orderID)
        if err != nil {
            http.Error(w, "Invalid order ID", http.StatusBadRequest)
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
            SELECT menu_item_id, SUM(quantity) AS total_quantity
            FROM order_items
            WHERE created_at BETWEEN $1 AND $2
            GROUP BY menu_item_id
        `
        rows, err := db.Query(query, startDate, endDate)
        if err != nil {
            http.Error(w, "Failed to fetch ordered items", http.StatusInternalServerError)
            return
        }
        defer rows.Close()

        result := make(map[string]int)
        for rows.Next() {
            var menuItemID string
            var totalQuantity int
            if err := rows.Scan(&menuItemID, &totalQuantity); err != nil {
                http.Error(w, "Failed to scan ordered items", http.StatusInternalServerError)
                return
            }
            result[menuItemID] = totalQuantity
        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(result)
    }
}