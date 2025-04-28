package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// FullTextSearchReport handles search across orders, menu items, and customers
func FullTextSearchReport(dbc *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		query := r.URL.Query().Get("q")
		if query == "" {
			http.Error(w, "Search query parameter 'q' is required", http.StatusBadRequest)
			return
		}

		filter := r.URL.Query().Get("filter")
		if filter == "" {
			filter = "all"
		}

		minPrice, _ := strconv.ParseFloat(r.URL.Query().Get("minPrice"), 64)
		maxPrice, _ := strconv.ParseFloat(r.URL.Query().Get("maxPrice"), 64)

		response := struct {
			MenuItems   []map[string]interface{} `json:"menu_items"`
			Orders      []map[string]interface{} `json:"orders"`
			TotalMatches int                     `json:"total_matches"`
		}{}

		// Search menu items if requested
		if filter == "all" || strings.Contains(filter, "menu") {
			menuQuery := `
				SELECT id, name, description, price,
					(CASE WHEN name ILIKE '%' || $1 || '%' THEN 1 ELSE 0 END) +
					(CASE WHEN description ILIKE '%' || $1 || '%' THEN 1 ELSE 0 END) AS relevance
				FROM menu_items
				WHERE name ILIKE '%' || $1 || '%' OR description ILIKE '%' || $1 || '%'
			`

			// Add price filtering if needed
			if minPrice > 0 || maxPrice > 0 {
				menuQuery += " AND price BETWEEN $2 AND $3"
			}

			// Add ordering and limiting (only once)
			menuQuery += " ORDER BY relevance DESC, name ASC LIMIT 10"

			var rows *sql.Rows
			var err error

			if minPrice > 0 || maxPrice > 0 {
				rows, err = dbc.Query(menuQuery, query, minPrice, maxPrice)
			} else {
				rows, err = dbc.Query(menuQuery, query)
			}

			if err != nil {
				http.Error(w, "Failed to search menu items: "+err.Error(), http.StatusInternalServerError)
				return
			}
			defer rows.Close()

			for rows.Next() {
				var item map[string]interface{} = make(map[string]interface{})
				var id, name, description string
				var price float64
				var relevance float32
				if err := rows.Scan(&id, &name, &description, &price, &relevance); err != nil {
					log.Println("Error scanning menu item:", err)
					continue
				}
				item["id"] = id
				item["name"] = name
				item["description"] = description
				item["price"] = price
				item["relevance"] = relevance
				response.MenuItems = append(response.MenuItems, item)
				response.TotalMatches++
			}
		}

		// Search orders if requested
		if filter == "all" || strings.Contains(filter, "orders") {
			orderQuery := `
				SELECT o.id, c.name, 
					array_agg(mi.name) as items, 
					o.total_amount,
					(CASE WHEN c.name ILIKE '%' || $1 || '%' THEN 1 ELSE 0 END) AS relevance
				FROM orders o
				JOIN customers c ON o.customer_id = c.id
				JOIN order_items oi ON o.id = oi.order_id
				JOIN menu_items mi ON oi.menu_item_id = mi.id
				WHERE c.name ILIKE '%' || $1 || '%'
			`
			
			// Add price filtering if needed
			if minPrice > 0 || maxPrice > 0 {
				orderQuery += " AND o.total_amount BETWEEN $2 AND $3"
			}
			
			// Add grouping and ordering (only once, at the end)
			orderQuery += " GROUP BY o.id, c.name, o.total_amount ORDER BY relevance DESC, o.id ASC LIMIT 10"
		
			var rows *sql.Rows
			var err error
			
			if minPrice > 0 || maxPrice > 0 {
				rows, err = dbc.Query(orderQuery, query, minPrice, maxPrice)
			} else {
				rows, err = dbc.Query(orderQuery, query)
			}

			if err != nil {
				http.Error(w, "Failed to search orders: "+err.Error(), http.StatusInternalServerError)
				return
			}
			defer rows.Close()

			for rows.Next() {
				var order map[string]interface{} = make(map[string]interface{})
				var id int
				var customerName string
				var items []string
				var total float64
				var relevance float32
				if err := rows.Scan(&id, &customerName, &items, &total, &relevance); err != nil {
					log.Println("Error scanning order:", err)
					continue
				}
				order["id"] = id
				order["customer_name"] = customerName
				order["items"] = items
				order["total"] = total
				order["relevance"] = relevance
				response.Orders = append(response.Orders, order)
				response.TotalMatches++
			}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

// OrderedItemsByPeriod returns order counts grouped by day or month
func OrderedItemsByPeriod(dbc *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		period := r.URL.Query().Get("period")
		if period != "day" && period != "month" {
			http.Error(w, "Invalid period parameter. Must be 'day' or 'month'", http.StatusBadRequest)
			return
		}

		response := struct {
			Period       string          `json:"period"`
			Month       string          `json:"month,omitempty"`
			Year        string          `json:"year,omitempty"`
			OrderedItems []map[string]int `json:"orderedItems"`
		}{
			Period: period,
		}

		if period == "day" {
			month := r.URL.Query().Get("month")
			if month == "" {
				month = time.Now().Month().String()
			}
			response.Month = strings.ToLower(month)

			// Get days in month
			year := time.Now().Year()
			monthNum := getMonthNumber(response.Month)
			daysInMonth := time.Date(year, time.Month(monthNum+1), 0, 0, 0, 0, 0, time.UTC).Day()

			query := `
				SELECT EXTRACT(DAY FROM created_at)::int as day, COUNT(*) as count
				FROM orders
				WHERE EXTRACT(MONTH FROM created_at) = $1
				GROUP BY day
				ORDER BY day
			`

			rows, err := dbc.Query(query, monthNum)
			if err != nil {
				http.Error(w, "Failed to query daily orders: "+err.Error(), http.StatusInternalServerError)
				return
			}
			defer rows.Close()

			// Initialize with all days set to 0
			dayCounts := make(map[int]int)
			for day := 1; day <= daysInMonth; day++ {
				dayCounts[day] = 0
			}

			// Update with actual counts from DB
			for rows.Next() {
				var day, count int
				if err := rows.Scan(&day, &count); err != nil {
					log.Println("Error scanning day count:", err)
					continue
				}
				dayCounts[day] = count
			}

			// Convert to required format
			for day := 1; day <= daysInMonth; day++ {
				response.OrderedItems = append(response.OrderedItems, map[string]int{
					strconv.Itoa(day): dayCounts[day],
				})
			}

		} else { // month
			year := r.URL.Query().Get("year")
			if year == "" {
				year = strconv.Itoa(time.Now().Year())
			}
			response.Year = year

			query := `
				SELECT EXTRACT(MONTH FROM created_at)::int as month, COUNT(*) as count
				FROM orders
				WHERE EXTRACT(YEAR FROM created_at) = $1
				GROUP BY month
				ORDER BY month
			`

			rows, err := dbc.Query(query, year)
			if err != nil {
				http.Error(w, "Failed to query monthly orders: "+err.Error(), http.StatusInternalServerError)
				return
			}
			defer rows.Close()

			// Initialize with all months set to 0
			monthCounts := make(map[int]int)
			for month := 1; month <= 12; month++ {
				monthCounts[month] = 0
			}

			// Update with actual counts from DB
			for rows.Next() {
				var month, count int
				if err := rows.Scan(&month, &count); err != nil {
					log.Println("Error scanning month count:", err)
					continue
				}
				monthCounts[month] = count
			}

			// Convert to required format
			for month := 1; month <= 12; month++ {
				response.OrderedItems = append(response.OrderedItems, map[string]int{
					time.Month(month).String(): monthCounts[month],
				})
			}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

// GetLeftovers returns inventory items with pagination and sorting
func GetLeftovers(dbc *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Parse query parameters
		sortBy := r.URL.Query().Get("sortBy")
		if sortBy != "price" && sortBy != "quantity" && sortBy != "" {
			http.Error(w, "Invalid sortBy parameter. Must be 'price' or 'quantity'", http.StatusBadRequest)
			return
		}

		page, err := strconv.Atoi(r.URL.Query().Get("page"))
		if err != nil || page < 1 {
			page = 1
		}

		pageSize, err := strconv.Atoi(r.URL.Query().Get("pageSize"))
		if err != nil || pageSize < 1 {
			pageSize = 10
		}

		// Build base query
		query := "SELECT name, stock as quantity, price FROM inventory"
		
		// Add sorting
		switch sortBy {
		case "price":
			query += " ORDER BY price DESC"
		case "quantity":
			query += " ORDER BY stock DESC"
		default:
			query += " ORDER BY name"
		}

		// Add pagination
		offset := (page - 1) * pageSize
		query += fmt.Sprintf(" LIMIT %d OFFSET %d", pageSize, offset)

		// Execute query
		rows, err := dbc.Query(query)
		if err != nil {
			http.Error(w, "Failed to query inventory: "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var items []map[string]interface{}
		for rows.Next() {
			var name string
			var quantity float64
			var price float64
			if err := rows.Scan(&name, &quantity, &price); err != nil {
				log.Println("Error scanning inventory item:", err)
				continue
			}
			items = append(items, map[string]interface{}{
				"name":     name,
				"quantity": quantity,
				"price":    price,
			})
		}

		// Get total count for pagination info
		var totalCount int
		err = dbc.QueryRow("SELECT COUNT(*) FROM inventory").Scan(&totalCount)
		if err != nil {
			http.Error(w, "Failed to count inventory items: "+err.Error(), http.StatusInternalServerError)
			return
		}

		totalPages := totalCount / pageSize
		if totalCount%pageSize != 0 {
			totalPages++
		}

		response := struct {
			CurrentPage int                      `json:"currentPage"`
			HasNextPage bool                     `json:"hasNextPage"`
			PageSize    int                      `json:"pageSize"`
			TotalPages  int                      `json:"totalPages"`
			Data        []map[string]interface{} `json:"data"`
		}{
			CurrentPage: page,
			HasNextPage: page < totalPages,
			PageSize:    pageSize,
			TotalPages:  totalPages,
			Data:        items,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

// BulkOrderProcess handles processing multiple orders in a transaction
func BulkOrderProcess(dbc *sql.DB) http.HandlerFunc {
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

		// Parse request body
		var request struct {
			Orders []struct {
				CustomerName string `json:"customer_name"`
				Items        []struct {
					MenuItemID string `json:"menu_item_id"`
					Quantity   int    `json:"quantity"`
				} `json:"items"`
			} `json:"orders"`
		}

		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		// Begin transaction
		tx, err := dbc.Begin()
		if err != nil {
			http.Error(w, "Failed to start transaction: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Prepare response structure
		response := struct {
			ProcessedOrders []map[string]interface{} `json:"processed_orders"`
			Summary         map[string]interface{}   `json:"summary"`
		}{
			ProcessedOrders: make([]map[string]interface{}, 0),
			Summary: map[string]interface{}{
				"total_orders": len(request.Orders),
				"accepted":     0,
				"rejected":     0,
				"total_revenue": 0.0,
				"inventory_updates": make([]map[string]interface{}, 0),
			},
		}

		// Process each order
		for _, order := range request.Orders {
			// Check inventory for all items first
			canProcess := true
			var totalAmount float64
			inventoryUpdates := make(map[string]int) // ingredientID -> quantity needed

			for _, item := range order.Items {
				// Get menu item ingredients and calculate needed quantities
				rows, err := tx.Query(`
					SELECT ingredient_id, quantity 
					FROM menu_item_ingredients 
					WHERE menu_item_id = $1`, item.MenuItemID)
				if err != nil {
					canProcess = false
					break
				}

				for rows.Next() {
					var ingredientID string
					var quantityPerUnit float64
					if err := rows.Scan(&ingredientID, &quantityPerUnit); err != nil {
						canProcess = false
						break
					}
					needed := quantityPerUnit * float64(item.Quantity)
					inventoryUpdates[ingredientID] += int(needed)
				}
				rows.Close()

				// Get menu item price
				var price float64
				err = tx.QueryRow("SELECT price FROM menu_items WHERE id = $1", item.MenuItemID).Scan(&price)
				if err != nil {
					canProcess = false
					break
				}
				totalAmount += price * float64(item.Quantity)
			}

			// Check inventory levels
			if canProcess {
				for ingredientID, needed := range inventoryUpdates {
					var currentStock float64
					err := tx.QueryRow("SELECT stock FROM inventory WHERE id = $1", ingredientID).Scan(&currentStock)
					if err != nil || currentStock < float64(needed) {
						canProcess = false
						break
					}
				}
			}
			var orderID int
			// Process order if possible
			if canProcess {
				// Create customer if not exists
				var customerID int
				err := tx.QueryRow(`
					INSERT INTO customers (name) 
					VALUES ($1) 
					ON CONFLICT (name) DO UPDATE SET name=EXCLUDED.name 
					RETURNING id`, order.CustomerName).Scan(&customerID)
				if err != nil {
					canProcess = false
				}

				// Create order
				
				err = tx.QueryRow(`
					INSERT INTO orders (customer_id, total_amount, status, payment_method) 
					VALUES ($1, $2, 'open', 'cash') 
					RETURNING id`, customerID, totalAmount).Scan(&orderID)
				if err != nil {
					canProcess = false
				}
				
				// Add order items
				if canProcess {
					for _, item := range order.Items {
						_, err := tx.Exec(`
							INSERT INTO order_items (order_id, menu_item_id, quantity, price_at_order)
							VALUES ($1, $2, $3, 
								(SELECT price FROM menu_items WHERE id = $2))`,
							orderID, item.MenuItemID, item.Quantity)
						if err != nil {
							canProcess = false
							break
						}
					}
				}

				// Update inventory if everything succeeded
				if canProcess {
					for ingredientID, needed := range inventoryUpdates {
						_, err := tx.Exec(`
							UPDATE inventory SET stock = stock - $1 
							WHERE id = $2`, needed, ingredientID)
						if err != nil {
							canProcess = false
							break
						}

						// Record inventory transaction
						_, err = tx.Exec(`
							INSERT INTO inventory_transactions 
							(inventory_id, change_amount, transaction_type) 
							VALUES ($1, $2, 'sale')`, ingredientID, -needed)
						if err != nil {
							canProcess = false
							break
						}
					}
				}
			}

			// Build response for this order
			orderResult := map[string]interface{}{
				"customer_name": order.CustomerName,
				"status":       "accepted",
				"total":        totalAmount,
			}
			
			if canProcess {
				orderResult["order_id"] = orderID
				response.Summary["accepted"] = response.Summary["accepted"].(int) + 1
				response.Summary["total_revenue"] = response.Summary["total_revenue"].(float64) + totalAmount
			} else {
				orderResult["status"] = "rejected"
				orderResult["reason"] = "insufficient_inventory"
				response.Summary["rejected"] = response.Summary["rejected"].(int) + 1
				// Rollback this order's changes
				tx.Rollback()
				// Start new transaction for remaining orders
				tx, err = dbc.Begin()
				if err != nil {
					http.Error(w, "Failed to restart transaction: "+err.Error(), http.StatusInternalServerError)
					return
				}
			}

			response.ProcessedOrders = append(response.ProcessedOrders, orderResult)
		}

		// Commit transaction if all succeeded
		if err := tx.Commit(); err != nil {
			http.Error(w, "Failed to commit transaction: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}
}

// Helper function to convert month name to number
func getMonthNumber(month string) int {
	months := map[string]int{
		"january":   1,
		"february":  2,
		"march":     3,
		"april":     4,
		"may":       5,
		"june":      6,
		"july":      7,
		"august":    8,
		"september": 9,
		"october":   10,
		"november":  11,
		"december":  12,
	}
	return months[strings.ToLower(month)]
}

func TotalAmount(dbc *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Invalid Method Request", http.StatusMethodNotAllowed)
			return
		}
		query := `Select SUM(total_amount) from orders`
		row := dbc.QueryRow(query)

		var sum float64
		if err := row.Scan(&sum); err != nil {
			http.Error(w, "Failed to scan ordered items", http.StatusInternalServerError)
			log.Println(err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(sum)
	}
}

func PopularItems(dbc *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Invalid Method Request", http.StatusMethodNotAllowed)
			return
		}
		query := `SELECT 
					mi.id AS menuItemID,
					mi.name,
					mi.description
					
				FROM 
					order_items oi
				JOIN 
					menu_items mi ON oi.menu_item_id = mi.id
				GROUP BY 
					mi.id, mi.name, mi.description`
		rows, err := dbc.Query(query)
		if err != nil {
			http.Error(w, "Failed to fetch ordered items", http.StatusInternalServerError)
			log.Println(err)
			return
		}
		defer rows.Close()
		type MenuItem struct {
			ID          string `json:"menuItemID"`
			Name        string `json:"name"`
			Description string `json:"description"`
		}
		var items []MenuItem
		for rows.Next() {
			var menuItemID string
			var name string
			var description string

			if err := rows.Scan(&menuItemID, &name, &description); err != nil {
				http.Error(w, "Failed to scan ordered items", http.StatusInternalServerError)
				log.Println(err)
				return
			}
			item := MenuItem{
				ID:          menuItemID,
				Name:        name,
				Description: description,
			}
			items = append(items, item)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(items)
	}
}
